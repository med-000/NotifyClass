package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

func (r *CourseRepository) Save(c *db.Course) error {
	var existing db.Course

	err := r.db.Where("external_id = ?", c.ExternalID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {

		r.log.Info.Printf("Create course external_id=%s", c.ExternalID)
		return r.db.Create(&c).Error

	} else if err != nil {
		r.log.Error.Printf("Save Error:,%v", err)
		return err
	}

	// update
	existing.Term = c.Term
	existing.Year = c.Year

	r.log.Info.Printf("Update couse existing_id=%d", c.ExternalID)
	return r.db.Save(&existing).Error

}
