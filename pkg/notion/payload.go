package notion

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/med-000/notifyclass/db"
)

func BuildClassPayload(databaseID string,c db.Class) map[string]interface{} {
	day:= dayToString(c.Day)

	year, term, err := parseCourseID(c.CourseID)
	if err != nil {
		log.Println(err)
		return nil
	}

	return map[string]interface{}{
		"parent": map[string]interface{}{
			"database_id": databaseID,
		},
		"properties": map[string]interface{}{
			"講義名": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": c.Title,
						},
					},
				},
			},
			"曜日": map[string]interface{}{
				"select": map[string]interface{}{
					"name": day,
				},
			},
			"時限": map[string]interface{}{
				"select": map[string]interface{}{
					"name": periodToString(c.Period),
				},
			},
			"開講年度": map[string]interface{}{
				"number": year,
			},
			"開講学期": map[string]interface{}{
				"number": term,
			},
		},
	}
}

func dayToString(day int) string {
	switch day {
	case 1:
		return "月曜"
	case 2:
		return "火曜"
	case 3:
		return "水曜"
	case 4:
		return "木曜"
	case 5:
		return "金曜"
	case 6:
		return "土曜"
	default:
		return ""
	}
}

func periodToString(p int) string {
	return fmt.Sprintf("%d限", p)
}

func parseCourseID(courseID string) (int, int, error) {
	parts := strings.Split(courseID, "_")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid course_id: %s", courseID)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year in course_id: %s", courseID)
	}

	term, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid term in course_id: %s", courseID)
	}

	return year, term, nil
}

func BuildEventPayload(databaseID string, e EventWithRelation) map[string]interface{} {

	var date map[string]interface{}

	if e.StartAt != nil {
		date = map[string]interface{}{
			"start": e.StartAt.Format(time.RFC3339),
		}
	}

	return map[string]interface{}{
		"parent": map[string]interface{}{
			"database_id": databaseID,
		},
		"properties": map[string]interface{}{
			"資料名": map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": e.Name,
						},
					},
				},
			},
			"グループ": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"content": e.Group,
						},
					},
				},
			},
			"種類": map[string]interface{}{
				"select": map[string]interface{}{
					"name": e.Category,
				},
			},
			"日付": map[string]interface{}{
				"date": date,
			},
			"講義名": map[string]interface{}{
				"relation": []map[string]interface{}{
					{
						"id": e.ClassPageID,
					},
				},
			},
		},
	}
}