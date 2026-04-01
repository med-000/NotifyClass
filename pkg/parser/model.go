package parser

type Course struct {
	Id     string
	Day    int
	Period int
	Title  string
	URL    string
}

type Class struct {
	Title  string
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
