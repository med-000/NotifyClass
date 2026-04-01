package service

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/scraping"
)

func FetchAll(req GetCourseRequest) (*parser.Course, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("els.sa.dendai.ac.jp"),
	)

	html, err := scraping.FetchCourseHTML(c, req.UserID, req.Password, req.Year, req.Term)
	if err != nil {
		log.Printf("ERROR: fetch failed: %v", err)
		return nil, err
	}

	classes := parser.ParseCourses(html)

	courseID := makeCourseID(req.Year, req.Term)

	var classDTOs []*parser.Class

	for i := range classes {
		classhtml, err := scraping.FetchClassHTML(c, classes[i].URL)
		if err != nil {
			log.Printf("WARN: class fetch failed: %v", err)
			continue
		}

		class := parser.ParseClass(classhtml)
		if class == nil {
			continue
		}

		classDTOs = append(classDTOs, &parser.Class{
			Id:     classes[i].Id,
			Day:    classes[i].Day,
			Period: classes[i].Period,
			Title:  classes[i].Title,
			URL:    classes[i].URL,
			Groups: class.Groups,
		})
	}

	return &parser.Course{
		Id:      courseID,
		Year:    req.Year,
		Term:    req.Term,
		Classes: classDTOs,
	}, nil
}

func makeCourseID(year int, term int) string {
	return fmt.Sprintf("%d_%d", year, term)
}
