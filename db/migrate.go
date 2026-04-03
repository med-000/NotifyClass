package db

import (
	"time"

	"gorm.io/gorm"
)

type MappingType string

const (
	MappingTypeClass MappingType = "class"
	MappingTypeEvent MappingType = "event"
)

type Course struct {
	ID string `gorm:"primaryKey"`

	Year int `gorm:"not null"`
	Term int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Class struct {
	ID uint `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null"`

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`
	
	Title  string
	Day    int
	Period int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Group struct {
	ID uint `gorm:"primaryKey"`

	ClassID uint
	Class Class `gorm:"foeignKey:ClassID"`

	Title string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID uint `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null"`

	GroupID uint
	Group Class `gorm:"foreignKey:GroupID"`

	Name string
	Category string
	StartAt *time.Time
	EndAt   *time.Time
	IsDone   bool `gorm:"default:false"`
	Notified bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotionMapping struct {
	ID uint `gorm:"primaryKey"`

	Type       MappingType `gorm:"type:varchar(20);not null;uniqueIndex:idx_external_type"`
	ExternalID string      `gorm:"type:varchar(128);not null;uniqueIndex:idx_external_type"`
	NotionPageID string `gorm:"type:varchar(128);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Course{},
		&Class{},
		&Event{},
		&NotionMapping{},
	)
}
