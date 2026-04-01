package db

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID uint `gorm:"primaryKey"`
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

	// 必須外部キー
	ClassID uint `gorm:"not null;index"`

	// 一意判定の軸
	Name string `gorm:"type:varchar(255);not null;uniqueIndex:idx_event_unique"`

	// 表示・分類
	Group    string `gorm:"type:varchar(100);not null;index"`
	Category string `gorm:"type:varchar(50);not null;index"`

	// 期間
	StartAt *time.Time `gorm:"uniqueIndex:idx_event_unique"`
	EndAt   *time.Time `gorm:"uniqueIndex:idx_event_unique"`

	// 状態
	IsDone   bool `gorm:"not null;default:false"`
	Notified bool `gorm:"not null;default:false"`

	// 外部連携
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
