package scraping

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func FetchClass(c *colly.Collector, url string) (string, error) {
	baseURL := "https://els.sa.dendai.ac.jp"
	redirectRe := regexp.MustCompile(`window\.location\.href\s*=\s*"([^"]+)"`)

	var html string

	cc := c.Clone()

	cc.OnResponse(func(r *colly.Response) {
		body := string(r.Body)

		// JSリダイレクト
		match := redirectRe.FindStringSubmatch(body)
		if len(match) > 1 {
			next := baseURL + match[1]
			_ = r.Request.Visit(next)
			return
		}

		// 最終ページ
		if strings.Contains(r.Request.URL.Path, "course.php") {
			html = body
		}
	})

	if err := cc.Visit(url); err != nil {
		return "", err
	}

	return html, nil
}
