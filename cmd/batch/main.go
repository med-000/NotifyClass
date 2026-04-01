package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/repository"
	"github.com/med-000/notifyclass/pkg/service"
)

func main() {
	// lock
	lockFile := "/tmp/notifyclass.lock"

	if _, err := os.Stat(lockFile); err == nil {
		log.Println("already running, skip")
		return
	}

	if err := os.WriteFile(lockFile, []byte("lock"), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(lockFile)

	// env
	_ = godotenv.Load()

	// DB
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	log.Println("DB ready")

	// fetch
	req := service.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     1,
	}

	courses, err := service.FetchCourses(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fetched courses: %d\n", len(courses))

	// save
	if err := repository.SaveCourses(database, courses); err != nil {
		log.Fatal(err)
	}

	log.Println("saved to DB")

	log.Println("done")
}
