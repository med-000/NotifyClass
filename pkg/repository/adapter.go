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

func ToDBGroup(g *parser.Group, classID uint) *db.Group {
	return &db.Group{
		ExternalID: g.ExternalId,
		ClassID:    classID,
		Title:      g.Name,
	}
}

func ToDBEvent(e *parser.Event, groupID uint) *db.Event {
	start, end := parseDate(e.Date)

	return &db.Event{
		ExternalID: e.ExternalId,
		GroupID:    groupID,
		Name:       e.Name,
		Category:   e.Category,
		StartAt:    start,
		EndAt:      end,
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