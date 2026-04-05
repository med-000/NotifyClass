package repository

import (
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
)

func (r *ContentRepository) Save(c *parser.Content) error {

	newContent := db.Content{
		CotentType: c.ContentType,
		FileName:   c.FileName,
		URL:        c.URL,
	}

	r.log.Info.Printf("Create Content FileName=%s", c.FileName)
	return r.db.Create(&newContent).Error

}
