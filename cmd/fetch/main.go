package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	course, err := appflow.FetchCourse()
	if err != nil {
		log.Fatal("failed to fetch course:", err)
	}

	if err := appflow.ExportCourseToJSON("course.json", course); err != nil {
		log.Fatal("failed to export json:", err)
	}

	log.Println("exported to course.json")
}
