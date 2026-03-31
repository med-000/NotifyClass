package service

import (
	"log"

	"github.com/gocolly/colly"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/scraping"
)

func FetchCourses(req GetCourseRequest) ([]CourseDTO, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("els.sa.dendai.ac.jp"),
	)

	html, err := scraping.FetchAll(c, req.UserID, req.Password, req.Year, req.Term)
	if err != nil {
		log.Printf("ERROR: fetch failed: %v", err)
		return nil, err
	}

	courses := parser.ParseCourses(html)

	var res []CourseDTO

	for i := range courses {
		classhtml, err := scraping.FetchClass(c, courses[i].URL)
		if err != nil {
			log.Printf("WARN: class fetch failed: %v", err)
			continue
		}

		class := parser.ParseClass(classhtml)
		if class == nil {
			continue
		}

		res = append(res, CourseDTO{
			Day:    courses[i].Day,
			Period: courses[i].Period,
			Title:  courses[i].Title,
			URL:    courses[i].URL,
			Groups: class.Groups,
		})
	}

	return res, nil
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

	courses := parser.ParseCourses(html)

	//対象の授業探す
	var targetURL string
	for _, course := range courses {
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
