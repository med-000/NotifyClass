package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
	"gorm.io/gorm"
)

func (r *EventRepository) Save(e *parser.Event) error {
	var existing db.Event

	err := r.db.Where("external_id = ?", e.ExternalId).First(&e).Error
	start, end := parseDate(e.Date)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newEvent := db.Event{
			ExternalID: e.ExternalId,
			Name:       e.Name,
			Category:   e.Category,
			StartAt:    start,
			EndAt:      end,
		}

		r.log.Info.Printf("Create event external_id=%s", e.ExternalId)
		return r.db.Create(&newEvent).Error

	} else if err != nil {
		r.log.Error.Printf("Save Error:,%v", err)
		return err
	}

	// update
	existing.Name = e.Name
	existing.Category = e.Category
	existing.StartAt = start
	existing.EndAt = end

	r.log.Info.Printf("Update couse existing_id=%d", e.ExternalId)
	return r.db.Save(&existing).Error

}

func parseDate(dateStr string) (*time.Time, *time.Time) {
	if dateStr == "" {
		return nil, nil
	}

	layout := "2006/01/02 15:04"
	parts := strings.Split(dateStr, " - ")
	if len(parts) != 2 {
		return nil, nil
	}

	start, err1 := time.Parse(layout, parts[0])
	end, err2 := time.Parse(layout, parts[1])

	if err1 != nil || err2 != nil {
		return nil, nil
	}

	return &start, &end
}
