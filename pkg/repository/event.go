package repository

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
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

		var class db.Class
		err := dbConn.
			Where("day = ? AND period = ?", course.Day, course.Period).
			First(&class).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			class = db.Class{
				Day:    course.Day,
				Period: course.Period,
				Title:  course.Title,
			}
			if err := dbConn.Create(&class).Error; err != nil {
				return err
			}
		}

		for _, group := range course.Groups {
			for _, ev := range group.Events {

				start, end := parseDate(ev.Date)

				event := db.Event{
					ClassID:  class.ID,
					Name:     ev.Name,
					Group:    group.Name,
					Category: ev.Category,
					StartAt:  start,
					EndAt:    end,
				}

				err := dbConn.
					Where("name = ? AND start_at <=> ? AND end_at <=> ?",
						event.Name, event.StartAt, event.EndAt).
					FirstOrCreate(&event).Error

				if err != nil {
					log.Println("insert error:", err)
				}
			}
		}
	}
	return nil
}