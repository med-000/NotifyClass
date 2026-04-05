package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

// FindByExternalID ExternalIDからClassを見つける
func (r *ClassRepository) FindByExternalID(externalID string) (*db.Class, error) {
	var c db.Class

	err := r.db.Where("external_id = ?", externalID).First(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Warn.Printf("class is nil: id=%s", externalID)
		return nil, nil
	}

	if err != nil {
		r.log.Error.Printf("FindByExternalID Error \nError Deitail:%s", err)
		return nil, err
	}

	r.log.Info.Printf("Success Found Class!")
	return &c, nil
}

// saveする
func (r *ClassRepository) Save(c *db.Class) error {
	var existing db.Class

	err := r.db.Where("external_id = ? AND course_id = ?", c.ExternalID, c.CourseID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Info.Printf("Create class external_id=%s", c.ExternalID)
		return r.db.Create(c).Error
	}

	if err != nil {
		r.log.Error.Printf("Save Error: %v", err)
		return err
	}

	// update
	existing.Day = c.Day
	existing.Period = c.Period
	existing.Title = c.Title

	r.log.Info.Printf("Update class external_id=%s", c.ExternalID)
	return r.db.Save(&existing).Error
}
