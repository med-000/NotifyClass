package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

func (r *GroupRepository) Save(g *db.Group) error {
	var existing db.Group

	err := r.db.Where("external_id = ?", g.ExternalID).First(&g).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {

		r.log.Info.Printf("Create Group external_id=%s", g.ExternalID)
		return r.db.Create(&g).Error
	} else if err != nil {
		r.log.Error.Printf("Save Error:%v", err)
		return err
	}

	//update
	existing.Title = g.Title

	r.log.Info.Printf("Update Group existing_id=%s", g.ExternalID)
	return r.db.Save(&existing).Error
}
