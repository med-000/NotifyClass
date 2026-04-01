package parser

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseCourses(html string) []Course {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	var courses []Course

	doc.Find("#schedule-table tbody tr").Each(func(i int, tr *goquery.Selection) {
		periodText := strings.TrimSpace(tr.Find(".schedule-table-class_order").Text())
		periodText = strings.ReplaceAll(periodText, "限", "")

		period, err := strconv.Atoi(periodText)
		if err != nil || period <= 0 {
			return
		}

		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			if j == 0 {
				return
			}

			a := td.Find("a")
			if a.Length() == 0 {
				return
			}

			title := strings.TrimSpace(a.Text())
			if title == "" {
				return
			}

			link, exists := a.Attr("href")
			if !exists || link == "" {
				return
			}

			url := link
			if !strings.HasPrefix(link, "http") {
				url = "https://els.sa.dendai.ac.jp" + link
			}

			id := ""
			parts := strings.Split(url, "/course.php/")
			if len(parts) > 1 {
				idPart := parts[1]
				id = strings.Split(idPart, "/")[0]
			}

			courses = append(courses, Course{
				Id:     id,
				Day:    j,
				Period: period,
				Title:  title,
				URL:    url,
			})
		})
	})

	return courses
}
