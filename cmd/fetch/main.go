package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	course, err := appflow.FetchCourse()
	if err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("fetch error: %v", err))
		log.Fatal("failed to fetch course:", err)
	}

	if err := appflow.ExportCourseToJSON("course.json", course); err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("json export error: %v", err))
		log.Fatal("failed to export json:", err)
	}

	log.Println("exported to course.json")
}
