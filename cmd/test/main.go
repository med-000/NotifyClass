package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/scraping"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("なんかおかしい")
	}
	scraperLogger, err := logger.NewScraperLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to init scraper logger: %v", err))
	}

	scraper := scraping.NewScraper(scraperLogger)

	req := scraping.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     1,
	}
	course, err := scraper.FetchAll(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = exportToJSON("course.json", course)
	if err != nil {
		fmt.Println("failed to export json:", err)
		return
	}

	fmt.Println("exported to course.json")
}

func exportToJSON(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
