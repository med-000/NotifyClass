package notion

import (
	"fmt"
	"time"

	"github.com/med-000/notifyclass/db"
)

func BuildClassPayload(cfg Config, c db.Class) map[string]interface{} {
	day := dayToString(c.Day)

	return map[string]interface{}{
		"parent": map[string]interface{}{
			"database_id": cfg.ClassDatabaseID,
		},
		"properties": map[string]interface{}{
			cfg.ClassTitleProperty:  titleProperty(c.Title),
			cfg.ClassDayProperty:    selectProperty(day),
			cfg.ClassPeriodProperty: selectProperty(periodToString(c.Period)),
			cfg.ClassYearProperty:   numberProperty(c.Course.Year),
			cfg.ClassTermProperty:   numberProperty(c.Course.Term),
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

func BuildEventPayload(cfg Config, e db.Event, classPageID string) map[string]interface{} {
	return map[string]interface{}{
		"parent": map[string]interface{}{
			"database_id": cfg.EventDatabaseID,
		},
		"properties": map[string]interface{}{
			cfg.EventTitleProperty:    titleProperty(e.Name),
			cfg.EventGroupProperty:    richTextProperty(e.GroupName),
			cfg.EventDoneProperty:     checkboxProperty(e.IsDone),
			cfg.EventDateProperty:     dateProperty(e.StartAt, e.EndAt),
			cfg.EventCategoryProperty: selectProperty(e.Category),
			cfg.EventClassProperty:    relationProperty(classPageID),
		},
	}
}

func titleProperty(content string) map[string]interface{} {
	return map[string]interface{}{
		"title": []map[string]interface{}{
			{
				"text": map[string]interface{}{
					"content": content,
				},
			},
		},
	}
}

func richTextProperty(content string) map[string]interface{} {
	return map[string]interface{}{
		"rich_text": []map[string]interface{}{
			{
				"text": map[string]interface{}{
					"content": content,
				},
			},
		},
	}
}

func selectProperty(name string) map[string]interface{} {
	return map[string]interface{}{
		"select": map[string]interface{}{
			"name": name,
		},
	}
}

func numberProperty(value int) map[string]interface{} {
	return map[string]interface{}{
		"number": value,
	}
}

func checkboxProperty(value bool) map[string]interface{} {
	return map[string]interface{}{
		"checkbox": value,
	}
}

func dateProperty(startAt, endAt *time.Time) map[string]interface{} {
	if startAt == nil {
		return map[string]interface{}{
			"date": nil,
		}
	}

	date := map[string]interface{}{
		"start": startAt.Format(time.RFC3339),
	}

	if endAt != nil {
		date["end"] = endAt.Format(time.RFC3339)
	}

	return map[string]interface{}{
		"date": date,
	}
}

func relationProperty(pageID string) map[string]interface{} {
	return map[string]interface{}{
		"relation": []map[string]interface{}{
			{
				"id": pageID,
			},
		},
	}
}
