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

func FindMappingByPageID(dbConn *gorm.DB, pageID string, t db.MappingType) (*db.NotionMapping, error) {
	var m db.NotionMapping

	err := dbConn.
		Where("notion_page_id = ? AND type = ?", pageID, t).
		First(&m).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func LoadMappingsByType(dbConn *gorm.DB, t db.MappingType) ([]db.NotionMapping, error) {
	var mappings []db.NotionMapping

	if err := dbConn.Where("type = ?", t).Find(&mappings).Error; err != nil {
		return nil, err
	}

	return mappings, nil
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

	return UpsertNotionWithPageID(dbConn, externalID, t, payload, notionPageIDFromMapping(m))
}

func UpsertNotionWithPageID(
	dbConn *gorm.DB,
	externalID string,
	t db.MappingType,
	payload any,
	existingPageID string,
) (string, error) {
	if existingPageID != "" {
		err := UpdateNotion(existingPageID, payload)
		if err != nil {
			return "", fmt.Errorf("update notion error: %w", err)
		}
		if err := UpsertMapping(dbConn, externalID, t, existingPageID); err != nil {
			return "", fmt.Errorf("touch mapping error: %w", err)
		}
		return existingPageID, nil
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

func notionPageIDFromMapping(m *db.NotionMapping) string {
	if m == nil {
		return ""
	}
	return m.NotionPageID
}
