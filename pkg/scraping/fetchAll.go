package scraping

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func FetchAll(c *colly.Collector, userId string, pass string, year int, term int) (string, error) {
	var html string
	var loggedIn bool

	redirectRe := regexp.MustCompile(`window\.location\.href="([^"]+)"`)

	loginURL := "https://els.sa.dendai.ac.jp/webclass/login.php"
	baseURL := "https://els.sa.dendai.ac.jp"

	c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)

		// JSリダイレクト
		if strings.Contains(r.Request.URL.String(), "login.php") {
			match := redirectRe.FindStringSubmatch(body)
			if len(match) > 1 {
				next := baseURL + match[1]
				r.Request.Visit(next)
				return
			}
		}

		// ログイン成功
		if strings.Contains(r.Request.URL.String(), "acs_=") {
			loggedIn = true

			// 学期変更
			r.Request.Post(
				baseURL+"/webclass/index.php",
				map[string]string{
					"year":     strconv.Itoa(int(year)),
					"semester": strconv.Itoa(int(term)),
				},
			)
			return
		}

		// 最終HTML
		if strings.Contains(r.Request.URL.String(), "index.php") {
			html = body
		}
	})

	// STEP1
	if err := c.Visit(loginURL); err != nil {
		return "", err
	}

	// STEP2
	if err := c.Post(loginURL, map[string]string{
		"username": userId,
		"val":      pass,
	}); err != nil {
		return "", err
	}

	// Clear sensitive data from memory
	// Note: In Go, strings are immutable, so we can't directly clear the memory
	// However, by reassigning empty strings, the original content can be garbage collected
	defer func() {
		// These reassignments allow the original data to be garbage collected faster
		userId = ""
		pass = ""
	}()

	if !loggedIn || html == "" {
		return "", errors.New("failed to fetch html")
	}

	return html, nil
}
