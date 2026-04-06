package notion

import "os"

type Config struct {
	APIToken string

	ClassDatabaseID string
	EventDatabaseID string

	ClassTitleProperty      string
	ClassPeriodProperty     string
	ClassDayProperty        string
	ClassAssignmentProperty string
	ClassYearProperty       string
	ClassTermProperty       string

	EventTitleProperty    string
	EventGroupProperty    string
	EventDoneProperty     string
	EventDateProperty     string
	EventCategoryProperty string
	EventClassProperty    string
}

type databaseQueryRequest struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

type databaseQueryResponse struct {
	Results    []notionPage `json:"results"`
	HasMore    bool         `json:"has_more"`
	NextCursor *string      `json:"next_cursor"`
}

type notionPage struct {
	ID         string                    `json:"id"`
	Properties map[string]notionProperty `json:"properties"`
}

type notionProperty struct {
	Type     string                 `json:"type"`
	Checkbox *bool                  `json:"checkbox,omitempty"`
	Date     map[string]interface{} `json:"date,omitempty"`
}

func LoadConfigFromEnv() Config {
	return Config{
		APIToken: os.Getenv("NOTION_API_KEY"),

		ClassDatabaseID: os.Getenv("NOTION_CLASS_DATABASE_ID"),
		EventDatabaseID: os.Getenv("NOTION_EVENT_DATABASE_ID"),

		ClassTitleProperty:      envOrDefault("NOTION_CLASS_TITLE_PROPERTY", "講義名"),
		ClassPeriodProperty:     envOrDefault("NOTION_CLASS_PERIOD_PROPERTY", "時限"),
		ClassDayProperty:        envOrDefault("NOTION_CLASS_DAY_PROPERTY", "曜日"),
		ClassAssignmentProperty: envOrDefault("NOTION_CLASS_ASSIGNMENT_PROPERTY", "課題"),
		ClassYearProperty:       envOrDefault("NOTION_CLASS_YEAR_PROPERTY", "開講年度"),
		ClassTermProperty:       envOrDefault("NOTION_CLASS_TERM_PROPERTY", "開講学期"),

		EventTitleProperty:    envOrDefault("NOTION_EVENT_TITLE_PROPERTY", "資料名"),
		EventGroupProperty:    envOrDefault("NOTION_EVENT_GROUP_PROPERTY", "グループ"),
		EventDoneProperty:     envOrDefault("NOTION_EVENT_DONE_PROPERTY", "完了"),
		EventDateProperty:     envOrDefault("NOTION_EVENT_DATE_PROPERTY", "日付"),
		EventCategoryProperty: envOrDefault("NOTION_EVENT_CATEGORY_PROPERTY", "種類"),
		EventClassProperty:    envOrDefault("NOTION_EVENT_CLASS_PROPERTY", "講義名"),
	}
}

func (c Config) Validate() error {
	switch {
	case c.APIToken == "":
		return ErrMissingNotionAPIKey
	case c.ClassDatabaseID == "":
		return ErrMissingClassDatabaseID
	case c.EventDatabaseID == "":
		return ErrMissingEventDatabaseID
	default:
		return nil
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
