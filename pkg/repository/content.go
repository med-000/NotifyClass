package repository

import (
	"github.com/med-000/notifyclass/db"
)

func (r *ContentRepository) Save(c *db.Content) error {
	r.log.Info.Printf("Create Content FileName=%s", c.FileName)
	return r.db.Create(&c).Error

}
