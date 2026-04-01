package service

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/scraping"
)

func FetchCourses(req GetCourseRequest) (*CourseDTO, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("els.sa.dendai.ac.jp"),
	)

	html, err := scraping.FetchAll(c, req.UserID, req.Password, req.Year, req.Term)
	if err != nil {
		log.Printf("ERROR: fetch failed: %v", err)
		return nil, err
	}

	classes := parser.ParseCourses(html, req.Year, req.Term)

	courseID := makeCourseID(req.Year, req.Term)

	var classDTOs []ClassDTO

	for i := range classes {
		classhtml, err := scraping.FetchClass(c, classes[i].URL)
		if err != nil {
			log.Printf("WARN: class fetch failed: %v", err)
			continue
		}

		class := parser.ParseClass(classhtml)
		if class == nil {
			continue
		}

		classDTOs = append(classDTOs, ClassDTO{
			Id:     classes[i].Id,
			Day:    classes[i].Day,
			Period: classes[i].Period,
			Title:  classes[i].Title,
			URL:    classes[i].URL,
			Groups: class.Groups,
		})
	}

	return &CourseDTO{
		Id:      courseID,
		Year:    req.Year,
		Term:    req.Term,
		Classes: classDTOs,
	}, nil
}

func FetchClassByRequest(req GetClassRequest) (*ClassDTO, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("els.sa.dendai.ac.jp"),
	)

	//一覧取得
	html, err := scraping.FetchAll(c, req.UserID, req.Password, req.Year, req.Term)
	if err != nil {
		log.Printf("ERROR: fetch failed: %v", err)
		return nil, err
	}

	classes := parser.ParseCourses(html, req.Year, req.Term)

	//対象の授業探す
	var targetURL string
	for _, course := range classes {
		if course.Day == req.Day && course.Period == req.Period {
			targetURL = course.URL
			break
		}
	}

	if targetURL == "" {
		return nil, nil
	}

	//詳細取得
	classhtml, err := scraping.FetchClass(c, targetURL)
	if err != nil {
		log.Printf("ERROR: class fetch failed: %v", err)
		return nil, err
	}

	class := parser.ParseClass(classhtml)
	if class == nil {
		return nil, nil
	}

	//DTO化
	return &ClassDTO{
		Title:  class.Title,
		Groups: class.Groups,
	}, nil
}

func makeCourseID(year int, term int) string {
	return fmt.Sprintf("%d_%d", year, term)
}
