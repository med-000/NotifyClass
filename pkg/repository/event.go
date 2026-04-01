package repository

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
	"gorm.io/gorm"
)

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

func SaveClasses(dbConn *gorm.DB, classes []*parser.Class) error {
	for _, class := range classes {

		// Class
		var dbclass db.Class

		err := dbConn.
			Where("external_id = ?", class.Id).
			First(&dbclass).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbclass = db.Class{
				ExternalID: class.Id,
				Day:        class.Day,
				Period:     class.Period,
				Title:      class.Title,
			}
			if err := dbConn.Create(&dbclass).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Event
		for _, group := range class.Groups {
			for _, ev := range group.Events {
				if ev.Id == "" {
					log.Println("not found Id")
					continue
				}

				start, end := parseDate(ev.Date)

				event := db.Event{
					ClassID:    dbclass.ID,
					ExternalID: ev.Id,
					Name:       ev.Name,
					Group:      group.Name,
					Category:   ev.Category,
					StartAt:    start,
					EndAt:      end,
				}

				err := dbConn.
					Where("external_id = ?", ev.Id).
					FirstOrCreate(&event).Error

				if err != nil {
					log.Println("insert error:", err)
					continue
				}

				// 更新検知
				var existing db.Event
				if err := dbConn.Where("external_id = ?", ev.Id).First(&existing).Error; err == nil {

					if changed(existing, ev, start, end, group.Name) {

						update := map[string]interface{}{
							"name":     ev.Name,
							"group":    group.Name,
							"category": ev.Category,
							"start_at": start,
							"end_at":   end,
							"notified": false,
						}

						if err := dbConn.Model(&existing).Updates(update).Error; err != nil {
							log.Println("update error:", err)
						}
					}
				}
			}
		}
	}
	return nil
}

func SaveCourse(dbConn *gorm.DB, course *parser.Course) error {

	//Course
	c := db.Course{
		ID:   course.Id,
		Year: course.Year,
		Term: course.Term,
	}

	if err := dbConn.FirstOrCreate(&c, db.Course{ID: c.ID}).Error; err != nil {
		return err
	}

	//Class
	for _, class := range course.Classes {

		var existing db.Class

		err := dbConn.
			Where("external_id = ? AND course_id = ?", class.Id, course.Id).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newClass := db.Class{
				ExternalID: class.Id,
				CourseID:   course.Id,
				Day:        class.Day,
				Period:     class.Period,
				Title:      class.Title,
			}

			if err := dbConn.Create(&newClass).Error; err != nil {
				return err
			}

			existing = newClass

		} else if err != nil {
			return err
		}

		//Event
		for _, group := range class.Groups {
			for _, ev := range group.Events {

				if ev.Id == "" {
					continue
				}

				start, end := parseDate(ev.Date)

				var existingEvent db.Event

				err := dbConn.
					Where("external_id = ? AND class_id = ?", ev.Id, existing.ID).
					First(&existingEvent).Error

				if errors.Is(err, gorm.ErrRecordNotFound) {
					event := db.Event{
						ClassID:    existing.ID,
						ExternalID: ev.Id,
						Name:       ev.Name,
						Group:      group.Name,
						Category:   ev.Category,
						StartAt:    start,
						EndAt:      end,
					}

					if err := dbConn.Create(&event).Error; err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func changed(e db.Event, ev *parser.Event, start, end *time.Time, group string) bool {

	if e.Name != ev.Name {
		return true
	}
	if e.Group != group {
		return true
	}
	if e.Category != ev.Category {
		return true
	}

	if !timeEqual(e.StartAt, start) {
		return true
	}
	if !timeEqual(e.EndAt, end) {
		return true
	}

	return false
}

func timeEqual(a, b *time.Time) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(*b)
}
