package discord

import "time"

type PendingTask struct {
	ExternalID string
	Name       string
	Category   string
	ClassTitle string
	Year       int
	Term       int
	Day        int
	Period     int
	Deadline   *time.Time
}

type SlotQuery struct {
	Year   int
	Term   int
	Day    int
	Period int
}
