package parser

type Course struct {
	Id   string
	Year int
	Term int
}

type Class struct {
	Id     string
	Day    int
	Period int
	Title  string
	URL    string
	Groups []*Group
}

type Group struct {
	Name   string
	Events []*Event
}

type Event struct {
	Id       string
	Name     string
	Category string
	Date     string
}
