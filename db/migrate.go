package db

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"type:varchar(64);uniqueIndex"`

	Day    int
	Period int
	Title  string

	NotionPageID *string `gorm:"type:varchar(255);uniqueIndex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
type Event struct {
	ID uint `gorm:"primaryKey"`

	ClassID uint `gorm:"not null;index"`

	ExternalID string `gorm:"type:varchar(64);uniqueIndex"`
	Name       string

	Group    string
	Category string

	StartAt *time.Time
	EndAt   *time.Time

	IsDone   bool `gorm:"default:false"`
	Notified bool `gorm:"default:false"`

	NotionPageID *string `gorm:"type:varchar(255);uniqueIndex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Class{},
		&Event{},
	)
}
