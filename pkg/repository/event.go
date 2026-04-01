package repository

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/service"
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

func SaveCourses(dbConn *gorm.DB, courses []service.CourseDTO) error {
	for _, course := range courses {

		// =========================
		// Class（IDベース）
		// =========================
		var class db.Class

		err := dbConn.
			Where("external_id = ?", course.Id).
			First(&class).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			class = db.Class{
				ExternalID: course.Id,
				Day:        course.Day,
				Period:     course.Period,
				Title:      course.Title,
			}
			if err := dbConn.Create(&class).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// =========================
		// Event
		// =========================
		for _, group := range course.Groups {
			for _, ev := range group.Events {

				// IDないやつは無視（掲示板など）
				if ev.Id == "" {
					continue
				}

				start, end := parseDate(ev.Date)

				var existing db.Event

				err := dbConn.
					Where("external_id = ?", ev.Id).
					First(&existing).Error

				if errors.Is(err, gorm.ErrRecordNotFound) {

					// --- 新規 ---
					event := db.Event{
						ClassID:    class.ID,
						ExternalID: ev.Id,
						Name:       ev.Name,
						Group:      group.Name,
						Category:   ev.Category,
						StartAt:    start,
						EndAt:      end,
					}

					if err := dbConn.Create(&event).Error; err != nil {
						log.Println("insert error:", err)
					}

				} else if err == nil {

					// --- 更新検知 ---
					if changed(existing, ev, start, end, group.Name) {

						update := map[string]interface{}{
							"name":     ev.Name,
							"group":    group.Name,
							"category": ev.Category,
							"start_at": start,
							"end_at":   end,

							// ←これが本質
							"notified": false,
						}

						if err := dbConn.Model(&existing).Updates(update).Error; err != nil {
							log.Println("update error:", err)
						}
					}

				} else {
					return err
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
