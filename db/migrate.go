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

type ContentType string

const (
	ContentTypePDF  ContentType = "pdf"
	ContentTypeHTML ContentType = "html"
	ContentTypeForm ContentType = "form"
)

type Course struct {
	ID         uint   `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null;uniqueIndex"`

	Year int `gorm:"not null"`
	Term int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Classes []Class `gorm:"foreignKey:CourseID"`
}

type Class struct {
	ID         uint   `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null;uniqueIndex"`

	CourseID uint
	Course   Course `gorm:"foreignKey:CourseID"`

	Title  string
	Day    int
	Period int

	CreatedAt time.Time
	UpdatedAt time.Time

	Groups []Group `gorm:"foreignKey:ClassID"`
}

type Group struct {
	ID         uint   `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null;uniqueIndex"`

	ClassID uint
	Class   Class `gorm:"foreignKey:ClassID"`

	Title string

	CreatedAt time.Time
	UpdatedAt time.Time

	Events []Event `gorm:"foreignKey:GroupID"`
}

type Event struct {
	ID         uint   `gorm:"primaryKey"`
	ExternalID string `gorm:"type:varchar(64);not null;uniqueIndex"`

	GroupID uint
	Group   Group `gorm:"foreignKey:GroupID"`

	Name     string
	Category string
	StartAt  *time.Time
	EndAt    *time.Time
	IsDone   bool `gorm:"default:false"`
	Notified bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Contents []Content `gorm:"foreignKey:EventID"`
}

type Content struct {
	ID uint `gorm:"primaryKey"`

	EventID uint
	Event   Event `gorm:"foreignKey:EventID"`

	ContentType ContentType `gorm:"type:varchar(20)"`
	URL         string
	FileName    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotionMapping struct {
	ID uint `gorm:"primaryKey"`

	Type         MappingType `gorm:"type:varchar(20);not null;uniqueIndex:idx_external_type"`
	ExternalID   string      `gorm:"type:varchar(128);not null;uniqueIndex:idx_external_type"`
	NotionPageID string      `gorm:"type:varchar(128);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Course{},
		&Class{},
		&Group{},
		&Event{},
		&Content{},
		&NotionMapping{},
	)
}