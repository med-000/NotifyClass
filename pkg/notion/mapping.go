package notion

import (
	"fmt"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindMapping(dbConn *gorm.DB, externalID string, t db.MappingType) (*db.NotionMapping, error) {
	var m db.NotionMapping

	err := dbConn.
		Where("external_id = ? AND type = ?", externalID, t).
		First(&m).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func SaveMapping(dbConn *gorm.DB, externalID string, t db.MappingType, pageID string) error {
	m := db.NotionMapping{
		ExternalID:   externalID,
		Type:         t,
		NotionPageID: pageID,
	}
	return dbConn.Create(&m).Error
}

func UpsertMapping(dbConn *gorm.DB, externalID string, t db.MappingType, pageID string) error {
	m := db.NotionMapping{
		ExternalID:   externalID,
		Type:         t,
		NotionPageID: pageID,
	}

	return dbConn.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "external_id"},
			{Name: "type"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"notion_page_id",
			"updated_at",
		}),
	}).Create(&m).Error
}

func UpsertNotion(
	dbConn *gorm.DB,
	externalID string,
	t db.MappingType,
	payload any,
) (string, error) {

	m, err := FindMapping(dbConn, externalID, t)
	if err != nil {
		return "", err
	}

	if m != nil {
		err := UpdateNotion(m.NotionPageID, payload)
		if err != nil {
			return "", fmt.Errorf("update notion error: %w", err)
		}
		return m.NotionPageID, nil
	}

	pageID, err := CreateNotion(payload)
	if err != nil {
		return "", fmt.Errorf("create notion error: %w", err)
	}

	err = UpsertMapping(dbConn, externalID, t, pageID)
	if err != nil {
		return "", fmt.Errorf("save mapping error: %w", err)
	}

	return pageID, nil
}
