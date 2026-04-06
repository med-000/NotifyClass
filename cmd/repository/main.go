package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(dbConn); err != nil {
		log.Fatal(err)
	}

	course, err := appflow.FetchCourse()
	if err != nil {
		log.Fatal("failed to fetch course:", err)
	}

	if err := appflow.SaveCourseToDB(dbConn, course); err != nil {
		log.Fatal("failed to save database:", err)
	}

	log.Println("saved to DB")
}
