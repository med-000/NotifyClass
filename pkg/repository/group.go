package repository

import (
	"errors"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
	"gorm.io/gorm"
)

func (r *GroupRepository) Save(g *parser.Group) error {
	var existing db.Group

	err := r.db.Where("external_id = ?", g.ExternalId).First(&g).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newGroup := db.Group{
			ExternalID: g.ExternalId,
			Title:      g.Name,
		}

		r.log.Info.Printf("Create Group external_id=%s", g.ExternalId)
		return r.db.Create(&newGroup).Error
	} else if err != nil {
		r.log.Error.Printf("Save Error:%v", err)
		return err
	}

	//update
	existing.Title = g.Name

	r.log.Info.Printf("Update Group existing_id=%d", g.ExternalId)
	return r.db.Save(&existing).Error
}
