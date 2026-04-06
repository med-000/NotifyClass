package repository

import (
	"time"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

type DueEvent struct {
	ExternalID string
	Name       string
	Category   string
	ClassTitle string
	Deadline   time.Time
}

type PendingTaskRow struct {
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

func (r *EventRepository) FindPendingDueEventsByDate(date time.Time) ([]DueEvent, error) {
	var events []DueEvent

	result := r.db.
		Table("events").
		Select(`
			events.external_id,
			events.name,
			events.category,
			classes.title as class_title,
			COALESCE(events.end_at, events.start_at) as deadline
		`).
		Joins("JOIN classes ON classes.id = events.class_id").
		Where("events.is_done = ?", false).
		Where("DATE(COALESCE(events.end_at, events.start_at)) = ?", date.Format("2006-01-02")).
		Order("deadline ASC").
		Scan(&events)

	if result.Error != nil {
		r.log.Error.Printf("FindPendingDueEventsByDate Error date=%s detail=%v", date.Format("2006-01-02"), result.Error)
		return nil, result.Error
	}

	r.log.Info.Printf("Found pending due events date=%s count=%d", date.Format("2006-01-02"), len(events))
	return events, nil
}

func (r *EventRepository) FindPendingTasks() ([]PendingTaskRow, error) {
	var tasks []PendingTaskRow

	result := r.pendingTaskBaseQuery().
		Order("deadline IS NULL ASC").
		Order("deadline ASC").
		Order("courses.year ASC").
		Order("courses.term ASC").
		Order("classes.day ASC").
		Order("classes.period ASC").
		Scan(&tasks)

	if result.Error != nil {
		r.log.Error.Printf("FindPendingTasks Error detail=%v", result.Error)
		return nil, result.Error
	}

	r.log.Info.Printf("Found pending tasks count=%d", len(tasks))
	return tasks, nil
}

func (r *EventRepository) FindPendingTasksBySlot(year, term, day, period int) ([]PendingTaskRow, error) {
	var tasks []PendingTaskRow

	result := r.pendingTaskBaseQuery().
		Where("courses.year = ?", year).
		Where("courses.term = ?", term).
		Where("classes.day = ?", day).
		Where("classes.period = ?", period).
		Order("deadline IS NULL ASC").
		Order("deadline ASC").
		Scan(&tasks)

	if result.Error != nil {
		r.log.Error.Printf(
			"FindPendingTasksBySlot Error year=%d term=%d day=%d period=%d detail=%v",
			year, term, day, period, result.Error,
		)
		return nil, result.Error
	}

	r.log.Info.Printf(
		"Found pending tasks by slot year=%d term=%d day=%d period=%d count=%d",
		year, term, day, period, len(tasks),
	)
	return tasks, nil
}

func (r *EventRepository) UpdateDoneByExternalID(externalID string, isDone bool) error {
	result := r.db.
		Model(&db.Event{}).
		Where("external_id = ?", externalID).
		Update("is_done", isDone)

	if result.Error != nil {
		r.log.Error.Printf("UpdateDoneByExternalID Error external_id=%s detail=%v", externalID, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Warn.Printf("event is nil: external_id=%s", externalID)
		return nil
	}

	r.log.Info.Printf("Update event is_done external_id=%s is_done=%t", externalID, isDone)
	return nil
}

func (r *EventRepository) pendingTaskBaseQuery() *gorm.DB {
	return r.db.
		Table("events").
		Select(`
			events.external_id,
			events.name,
			events.category,
			classes.title as class_title,
			courses.year,
			courses.term,
			classes.day,
			classes.period,
			COALESCE(events.end_at, events.start_at) as deadline
		`).
		Joins("JOIN classes ON classes.id = events.class_id").
		Joins("JOIN courses ON courses.id = classes.course_id").
		Where("events.is_done = ?", false)
}
