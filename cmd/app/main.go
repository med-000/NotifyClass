package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	log.Println("DB ready")

	if err := appflow.SyncNotionPull(database); err != nil {
		log.Fatal("failed to sync notion pull:", err)
	}

	course, err := appflow.FetchCourse()
	if err != nil {
		log.Fatal("failed to fetch course:", err)
	}

	if err := appflow.ExportCourseToJSON("course.json", course); err != nil {
		log.Fatal("failed to export json:", err)
	}

	log.Println("exported to course.json")
	if err := appflow.SaveCourseToDB(database, course); err != nil {
		log.Fatal("failed to save database:", err)
	}

	if err := appflow.SyncNotionPush(database); err != nil {
		log.Fatal("failed to sync notion push:", err)
	}
}
