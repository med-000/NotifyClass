package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

func (r *ContentRepository) Save(c *db.Content) error {
	var existing db.Content

	err := r.db.Where("url = ?", c.URL).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Info.Printf("Create cotent url=%s", c.URL)
		return r.db.Create(c).Error

	} 
	if err != nil {
		r.log.Error.Printf("Save Error:,%v", err)
		return err
	}

	// update
	existing.ContentType = c.ContentType
	existing.URL = c.URL
	existing.FileName = c.FileName

	r.log.Info.Printf("Update content URL=%s", c.URL)
	return r.db.Save(&existing).Error
}
