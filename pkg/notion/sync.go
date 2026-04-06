package notion

import (
	"fmt"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/logger"
	"gorm.io/gorm"
)

func SyncAllFromDB(dbConn *gorm.DB, notionLog *logger.NotionLogger, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	classes, err := loadClasses(dbConn)
	if err != nil {
		return fmt.Errorf("load classes: %w", err)
	}

	classMappings, err := loadMappingByExternalID(dbConn, db.MappingTypeClass)
	if err != nil {
		return fmt.Errorf("load class mappings: %w", err)
	}

	if err := SyncClassesToNotion(dbConn, notionLog, classes, classMappings, cfg); err != nil {
		return fmt.Errorf("sync classes: %w", err)
	}

	events, err := loadEvents(dbConn)
	if err != nil {
		return fmt.Errorf("load events: %w", err)
	}

	eventMappings, err := loadMappingByExternalID(dbConn, db.MappingTypeEvent)
	if err != nil {
		return fmt.Errorf("load event mappings: %w", err)
	}

	if err := SyncEventsToNotion(dbConn, notionLog, events, classMappings, eventMappings, cfg); err != nil {
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
	eventMappingsByPageID, err := loadMappingByPageID(dbConn, db.MappingTypeEvent)
	if err != nil {
		return fmt.Errorf("load event mappings by page id: %w", err)
	}

	for {
		response, err := QueryDatabase(cfg.EventDatabaseID, databaseQueryRequest{
			StartCursor: startCursor,
			PageSize:    100,
		})
		if err != nil {
			return fmt.Errorf("query event database: %w", err)
		}

		for _, page := range response.Results {
			if err := syncSingleEventCompletion(dbConn, notionLog, cfg, page, eventMappingsByPageID); err != nil {
				notionLog.Error.Printf("sync event completion failed notion_page_id=%s err=%v", page.ID, err)
				continue
			}
		}

		if !response.HasMore || response.NextCursor == nil || *response.NextCursor == "" {
			break
		}

		startCursor = *response.NextCursor
	}

	notionLog.Info.Printf("event completion sync from notion done")

	return nil
}

func SyncClassesToNotion(
	dbConn *gorm.DB,
	notionLog *logger.NotionLogger,
	classes []db.Class,
	classMappings map[string]string,
	cfg Config,
) error {
	for _, class := range classes {
		payload := BuildClassPayload(cfg, class)

		pageID, err := UpsertNotionWithPageID(
			dbConn,
			class.ExternalID,
			db.MappingTypeClass,
			payload,
			classMappings[class.ExternalID],
		)
		if err != nil {
			notionLog.Error.Printf("class sync failed external_id=%s title=%s err=%v", class.ExternalID, class.Title, err)
			continue
		}

		classMappings[class.ExternalID] = pageID
		notionLog.Info.Printf("class synced external_id=%s notion_page_id=%s", class.ExternalID, pageID)
	}

	return nil
}

func SyncEventsToNotion(
	dbConn *gorm.DB,
	notionLog *logger.NotionLogger,
	events []db.Event,
	classMappings map[string]string,
	eventMappings map[string]string,
	cfg Config,
) error {
	for _, event := range events {
		classPageID := classMappings[event.Class.ExternalID]
		if classPageID == "" {
			notionLog.Warn.Printf("class mapping missing class_external_id=%s", event.Class.ExternalID)
			continue
		}

		payload := BuildEventPayload(cfg, event, classPageID)

		pageID, err := upsertWithRetry(
			dbConn,
			notionLog,
			event.ExternalID,
			db.MappingTypeEvent,
			payload,
			eventMappings[event.ExternalID],
			3,
		)
		if err != nil {
			notionLog.Error.Printf("event sync failed external_id=%s name=%s err=%v", event.ExternalID, event.Name, err)
			continue
		}

		eventMappings[event.ExternalID] = pageID
		notionLog.Info.Printf("event synced external_id=%s notion_page_id=%s", event.ExternalID, pageID)
	}

	return nil
}

func syncSingleEventCompletion(
	dbConn *gorm.DB,
	notionLog *logger.NotionLogger,
	cfg Config,
	page notionPage,
	eventMappingsByPageID map[string]string,
) error {
	externalID := eventMappingsByPageID[page.ID]
	if externalID == "" {
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
		Where("external_id = ?", externalID).
		Update("is_done", done)

	if result.Error != nil {
		return fmt.Errorf("update event is_done: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		notionLog.Warn.Printf("event not found external_id=%s notion_page_id=%s", externalID, page.ID)
		return nil
	}

	notionLog.Info.Printf("event completion synced external_id=%s notion_page_id=%s is_done=%t", externalID, page.ID, done)
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
	existingPageID string,
	maxAttempts int,
) (string, error) {
	var (
		pageID string
		err    error
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pageID, err = UpsertNotionWithPageID(dbConn, externalID, mappingType, payload, existingPageID)
		if err == nil {
			return pageID, nil
		}

		notionLog.Warn.Printf("notion retry attempt=%d external_id=%s err=%v", attempt, externalID, err)
	}

	return "", err
}

func loadMappingByExternalID(dbConn *gorm.DB, t db.MappingType) (map[string]string, error) {
	mappings, err := LoadMappingsByType(dbConn, t)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(mappings))
	for _, mapping := range mappings {
		result[mapping.ExternalID] = mapping.NotionPageID
	}

	return result, nil
}

func loadMappingByPageID(dbConn *gorm.DB, t db.MappingType) (map[string]string, error) {
	mappings, err := LoadMappingsByType(dbConn, t)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(mappings))
	for _, mapping := range mappings {
		result[mapping.NotionPageID] = mapping.ExternalID
	}

	return result, nil
}
