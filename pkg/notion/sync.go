package notion

import (
	"fmt"
	"time"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/logger"
	"gorm.io/gorm"
)

const notionRateLimitInterval = 350 * time.Millisecond

func SyncAllFromDB(dbConn *gorm.DB, notionLog *logger.NotionLogger, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	classes, err := loadClasses(dbConn)
	if err != nil {
		return fmt.Errorf("load classes: %w", err)
	}

	if err := SyncClassesToNotion(dbConn, notionLog, classes, cfg); err != nil {
		return fmt.Errorf("sync classes: %w", err)
	}

	events, err := loadEvents(dbConn)
	if err != nil {
		return fmt.Errorf("load events: %w", err)
	}

	if err := SyncEventsToNotion(dbConn, notionLog, events, cfg); err != nil {
		return fmt.Errorf("sync events: %w", err)
	}

	return nil
}

func SyncEventCompletionFromNotion(dbConn *gorm.DB, notionLog *logger.NotionLogger, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	notionLog.Info.Printf("start sync event completion from notion")

	startCursor := ""

	for {
		response, err := QueryDatabase(cfg.EventDatabaseID, databaseQueryRequest{
			StartCursor: startCursor,
			PageSize:    100,
		})
		if err != nil {
			return fmt.Errorf("query event database: %w", err)
		}

		for _, page := range response.Results {
			if err := syncSingleEventCompletion(dbConn, notionLog, cfg, page); err != nil {
				notionLog.Error.Printf("sync event completion failed notion_page_id=%s err=%v", page.ID, err)
				continue
			}
			time.Sleep(notionRateLimitInterval)
		}

		if !response.HasMore || response.NextCursor == nil || *response.NextCursor == "" {
			break
		}

		startCursor = *response.NextCursor
	}

	notionLog.Info.Printf("event completion sync from notion done")

	return nil
}

func SyncClassesToNotion(dbConn *gorm.DB, notionLog *logger.NotionLogger, classes []db.Class, cfg Config) error {
	for _, class := range classes {
		payload := BuildClassPayload(cfg, class)

		pageID, err := UpsertNotion(dbConn, class.ExternalID, db.MappingTypeClass, payload)
		if err != nil {
			notionLog.Error.Printf("class sync failed external_id=%s title=%s err=%v", class.ExternalID, class.Title, err)
			continue
		}

		notionLog.Info.Printf("class synced external_id=%s notion_page_id=%s", class.ExternalID, pageID)
		time.Sleep(notionRateLimitInterval)
	}

	return nil
}

func SyncEventsToNotion(dbConn *gorm.DB, notionLog *logger.NotionLogger, events []db.Event, cfg Config) error {
	for _, event := range events {
		classMapping, err := FindMapping(dbConn, event.Class.ExternalID, db.MappingTypeClass)
		if err != nil {
			notionLog.Error.Printf("class mapping lookup failed class_external_id=%s err=%v", event.Class.ExternalID, err)
			continue
		}
		if classMapping == nil {
			notionLog.Warn.Printf("class mapping missing class_external_id=%s", event.Class.ExternalID)
			continue
		}

		payload := BuildEventPayload(cfg, event, classMapping.NotionPageID)

		pageID, err := upsertWithRetry(dbConn, notionLog, event.ExternalID, db.MappingTypeEvent, payload, 3)
		if err != nil {
			notionLog.Error.Printf("event sync failed external_id=%s name=%s err=%v", event.ExternalID, event.Name, err)
			continue
		}

		notionLog.Info.Printf("event synced external_id=%s notion_page_id=%s", event.ExternalID, pageID)
		time.Sleep(notionRateLimitInterval)
	}

	return nil
}

func syncSingleEventCompletion(
	dbConn *gorm.DB,
	notionLog *logger.NotionLogger,
	cfg Config,
	page notionPage,
) error {
	mapping, err := FindMappingByPageID(dbConn, page.ID, db.MappingTypeEvent)
	if err != nil {
		return fmt.Errorf("find mapping by page id: %w", err)
	}
	if mapping == nil {
		notionLog.Warn.Printf("event mapping missing notion_page_id=%s", page.ID)
		return nil
	}

	property, ok := page.Properties[cfg.EventDoneProperty]
	if !ok {
		notionLog.Warn.Printf("done property missing notion_page_id=%s property=%s", page.ID, cfg.EventDoneProperty)
		return nil
	}

	if property.Type != "checkbox" || property.Checkbox == nil {
		notionLog.Warn.Printf("done property invalid notion_page_id=%s property=%s type=%s", page.ID, cfg.EventDoneProperty, property.Type)
		return nil
	}

	done := *property.Checkbox
	result := dbConn.
		Model(&db.Event{}).
		Where("external_id = ?", mapping.ExternalID).
		Update("is_done", done)

	if result.Error != nil {
		return fmt.Errorf("update event is_done: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		notionLog.Warn.Printf("event not found external_id=%s notion_page_id=%s", mapping.ExternalID, page.ID)
		return nil
	}

	notionLog.Info.Printf("event completion synced external_id=%s notion_page_id=%s is_done=%t", mapping.ExternalID, page.ID, done)
	return nil
}

func loadClasses(dbConn *gorm.DB) ([]db.Class, error) {
	var classes []db.Class

	if err := dbConn.Preload("Course").Find(&classes).Error; err != nil {
		return nil, err
	}

	return classes, nil
}

func loadEvents(dbConn *gorm.DB) ([]db.Event, error) {
	var events []db.Event

	if err := dbConn.Preload("Class").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func upsertWithRetry(
	dbConn *gorm.DB,
	notionLog *logger.NotionLogger,
	externalID string,
	mappingType db.MappingType,
	payload any,
	maxAttempts int,
) (string, error) {
	var (
		pageID string
		err    error
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pageID, err = UpsertNotion(dbConn, externalID, mappingType, payload)
		if err == nil {
			return pageID, nil
		}

		notionLog.Warn.Printf("notion retry attempt=%d external_id=%s err=%v", attempt, externalID, err)
		time.Sleep(1 * time.Second)
	}

	return "", err
}
