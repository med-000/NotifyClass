package service

import "github.com/med-000/notifyclass/pkg/parser"

type GetCourseRequest struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	Year     int    `json:"year"`
	Term     int    `json:"term"`
}

type CourseDTO struct {
	Id      string     `json:"id"`
	Year    int        `json:"day"`
	Term    int        `json:"period"`
	Classes []ClassDTO `json:"classes"`
}

type CourseResponse struct {
	Courses []CourseDTO `json:"courses"`
}

type GetClassRequest struct {
	Id       string `json:"id"`
	UserID   string `json:"userId"`
	Password string `json:"password"`
	Year     int    `json:"year"`
	Term     int    `json:"term"`
	Day      int    `json:"day"`
	Period   int    `json:"period"`
}

type ClassDTO struct {
	Id     string          `json:"id"`
	Day    int             `json:"day"`
	Period int             `json:"period"`
	Title  string          `json:"title"`
	URL    string          `json:"url"`
	Groups []*parser.Group `json:"groups"`
}

type ClassResponse struct {
	Class *ClassDTO `json:"class"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}
