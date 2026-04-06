package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	courses, err := appflow.FetchCourses(time.Now())
	if err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("fetch error: %v", err))
		log.Fatal("failed to fetch courses:", err)
	}

	if err := appflow.ExportCourseToJSON("course.json", courses); err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("json export error: %v", err))
		log.Fatal("failed to export json:", err)
	}

	log.Println("exported to course.json")
}
