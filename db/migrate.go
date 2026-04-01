package db

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID   string `gorm:"primaryKey"` // "2025_1"
	Year int    `gorm:"not null"`
	Term int    `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Class struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"type:varchar(64);not null"`

	CourseID string `gorm:"not null;index"`

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

	ExternalID string `gorm:"type:varchar(64);not null"`

	Name string

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
		&Course{},
		&Class{},
		&Event{},
	)
}
