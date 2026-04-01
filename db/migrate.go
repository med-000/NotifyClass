package db

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID uint `gorm:"primaryKey"`

	// coursesと紐づけ（将来用）
	CourseName string
	Day        int
	Period     int

	// 授業ページ情報
	Title string

	// Notion連携
	NotionPageID *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
type Event struct {
	ID uint `gorm:"primaryKey"`

	ClassID *uint

	CourseName string `gorm:"type:varchar(255)"`
	Day        int
	Period     int

	Name     string `gorm:"type:varchar(255);uniqueIndex:idx_event_unique"`
	Category string `gorm:"type:varchar(100)"`

	StartAt *time.Time `gorm:"uniqueIndex:idx_event_unique"`
	EndAt   *time.Time `gorm:"uniqueIndex:idx_event_unique"`

	IsDone   bool
	Notified bool

	NotionPageID *string `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Class{},
		&Event{},
	)
}
