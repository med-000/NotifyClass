package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
	"gorm.io/gorm"
)

func (r *CourseRepository) Save(c *parser.Course) error {
	var existing db.Course

	err := r.db.Where("external_id = ?", c.ExternalId).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newCourse := db.Course{
			ExternalID: c.ExternalId,
			Year:       c.Year,
			Term:       c.Term,
		}

		r.log.Info.Printf("Create course external_id=%s", c.ExternalId)
		return r.db.Create(&newCourse).Error

	} else if err != nil {
		r.log.Error.Printf("Save Error:,%v", err)
		return err
	}

	// update
	existing.Term = c.Term
	existing.Year = c.Year

	r.log.Info.Printf("Update couse existing_id=%d", c.ExternalId)
	return r.db.Save(&existing).Error

}
