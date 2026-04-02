package notion

type NotionPageRequest struct {
	Parent     Parent              `json:"parent"`
	Properties map[string]Property `json:"properties"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Property struct {
	Title  []RichText `json:"title,omitempty"`
	Select *Select    `json:"select,omitempty"`
	Number *int       `json:"number,omitempty"`
}

type RichText struct {
	Text Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type Select struct {
	Name string `json:"name"`
}
