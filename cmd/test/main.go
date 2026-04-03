package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/scraping"
)

func main() {
	_ = godotenv.Load()
	scraperLogger, _ := logger.NewScraperLogger()

	scraper := scraping.NewScraper(scraperLogger)

	req := scraping.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2026,
		Term:     1,
	}
	course, err := scraper.FetchAll(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(course)
}
