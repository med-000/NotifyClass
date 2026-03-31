package handler

import "github.com/med-000/notifyclass/pkg/parser"

type GetCourseRequest struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	Year     int16  `json:"year"`
	Term     int16  `json:"term"`
	Day      int    `json:"day"`
	Period   int    `json:"period"`
}

type CourseDTO struct {
	Day    int             `json:"day"`
	Period int             `json:"period"`
	Title  string          `json:"title"`
	URL    string          `json:"url"`
	Groups []*parser.Group `json:"groups"`
}

type CourseResponse struct {
	Courses []CourseDTO `json:"courses"`
}

type ClassDTO struct {
	Title  string          `json:"title"`
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
