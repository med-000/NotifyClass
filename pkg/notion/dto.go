package notion

import "time"

type EventWithRelation struct {
	ExternalID  string
	Name        string
	Group       string
	Category    string
	StartAt     *time.Time
	ClassPageID string
}

type SyncEvent struct {
	EventExternalID string `gorm:"column:event_external_id"`
	Name            string
	Group           string
	Category        string
	StartAt         *time.Time
	ClassExternalID string `gorm:"column:class_external_id"`
}
