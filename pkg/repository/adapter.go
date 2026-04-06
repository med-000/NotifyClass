package repository

import (
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
)

func ToDBCourse(c *parser.Course) *db.Course {
	return &db.Course{
		ExternalID: c.ExternalId,
		Year:       c.Year,
		Term:       c.Term,
	}
}

func ToDBClass(c *parser.Class, courseID uint) *db.Class {
	return &db.Class{
		ExternalID: c.ExternalId,
		CourseID:   courseID,
		Title:      c.Title,
		Day:        c.Day,
		Period:     c.Period,
	}
}

func ToDBEvent(e *parser.Event, classID uint) *db.Event {
	start, end := parseDate(e.Date)

	return &db.Event{
		ExternalID: e.ExternalId,
		ClassID:    classID,
		Name:       e.Name,
		Category:   e.Category,
		GroupName:  e.GroupName,
		StartAt:    start,
		EndAt:      end,
		IsDone:     defaultEventIsDone(e.Category),
	}
}

func ToDBContent(c *parser.Content, eventID uint) *db.Content {
	return &db.Content{
		EventID:     eventID,
		ContentType: c.ContentType,
		URL:         c.URL,
		FileName:    c.FileName,
	}
}

func parseDate(dateStr string) (*time.Time, *time.Time) {
	if dateStr == "" {
		return nil, nil
	}

	layout := "2006/01/02 15:04"
	parts := strings.Split(dateStr, " - ")
	if len(parts) != 2 {
		return nil, nil
	}

	start, err1 := time.Parse(layout, parts[0])
	end, err2 := time.Parse(layout, parts[1])

	if err1 != nil || err2 != nil {
		return nil, nil
	}

	return &start, &end
}

func defaultEventIsDone(category string) bool {
	switch strings.TrimSpace(category) {
	case "レポート", "試験", "レポート（成績非公開）":
		return false
	default:
		return true
	}
}
