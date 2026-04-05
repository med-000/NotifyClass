package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

func (r *EventRepository) Save(e *db.Event) error {
	var existing db.Event

	err := r.db.Where("external_id = ?", e.ExternalID).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {

		r.log.Info.Printf("Create event external_id=%s", e.ExternalID)
		return r.db.Create(&e).Error

	} else if err != nil {
		r.log.Error.Printf("Save Error:,%v", err)
		return err
	}

	// update
	existing.Name = e.Name
	existing.Category = e.Category
	existing.StartAt = e.StartAt
	existing.EndAt = e.EndAt

	r.log.Info.Printf("Update couse existing_id=%s", e.ExternalID)
	return r.db.Save(&existing).Error

}
