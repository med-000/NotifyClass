package parser

type Course struct {
	ExternalId string
	Year       int
	Term       int
	Classes    []*Class
}

type Class struct {
	ExternalId string
	Day        int
	Period     int
	Title      string
	URL        string
	Groups     []*Group
}

type Group struct {
	Name   string
	Events []*Event
}

type Event struct {
	ExternalId string
	Name       string
	Category   string
	Date       string
	URL        string
	Content    []*Content
}

type Content struct {
	Type     string
	URL      string
	FileName string
}
