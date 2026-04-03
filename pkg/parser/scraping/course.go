package scraping

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseCourses(html string) []*Class {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	var classes []*Class

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

			fullURL := link
			if !strings.HasPrefix(link, "http") {
				fullURL = "https://els.sa.dendai.ac.jp" + link
			}

			// --- ID抽出（安全版） ---
			id := ""
			if idx := strings.Index(fullURL, "/course.php/"); idx != -1 {
				idPart := fullURL[idx+len("/course.php/"):]
				id = strings.Split(idPart, "/")[0]
			}

			classes = append(classes, &Class{
				Id:     id,
				Day:    j,
				Period: period,
				Title:  title,
				URL:    fullURL,
			})
		})
	})

	return classes
}
