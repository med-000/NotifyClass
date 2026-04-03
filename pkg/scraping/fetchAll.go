package scraping

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/med-000/notifyclass/pkg/parser/scraping"
)

func (s *Scraper) FetchAll(req GetCourseRequest) (*scraping.Course, error) {
	allowdomain := os.Getenv("ALLOW_DOMAIN")

	//collyのinit
	c := colly.NewCollector(
		colly.AllowedDomains(allowdomain),
	)

	//受け取ったrequestからhtmlを取得
	html, err := s.FetchCourseHTML(c, req.UserID, req.Password, req.Year, req.Term)
	if err != nil {
		s.log.Error.Printf("Not Get CourseHTML! error detail:%d", err)
		return nil, err
	}
	s.log.Info.Printf("Success Get CourseHTML!")

	classes := scraping.ParseCourses(html)
	s.log.Info.Printf("Success Parser CourseHTML!")

	//courseIDのフォーマット変換
	courseID := makeCourseID(req.Year, req.Term)

	var result []*scraping.Class

	//classesを展開してその数だけ継続(course内にあった講義数分中のものを取得してappend)
	for i := range classes {
		classhtml, err := s.FetchClassHTML(c, classes[i].URL)
		if err != nil {
			s.log.Error.Printf("Failed FetchClassHTML! error detail:%d", err)
			continue
		}

		class := scraping.ParseClass(classhtml)
		if class == nil {
			s.log.Info.Printf("classhtml is nil")
			continue
		}

		result = append(result, &scraping.Class{
			Id:     classes[i].Id,
			Day:    classes[i].Day,
			Period: classes[i].Period,
			Title:  classes[i].Title,
			URL:    classes[i].URL,
			Groups: class.Groups,
		})
	}

	return &scraping.Course{
		Id:      courseID,
		Year:    req.Year,
		Term:    req.Term,
		Classes: result,
	}, nil
}

//CourseIDの変換関数
func makeCourseID(year int, term int) string {
	return fmt.Sprintf("%d_%d", year, term)
}
