package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	myyearStr := "20" + os.Getenv("USER_ID")[:2]

	myyear, err := strconv.Atoi(myyearStr)
	if err != nil {
		return
	}
	year := time.Now().Year()

	for y := myyear; y <= year; y++ {
		for term := 1; term <= 2; term++ {
			fmt.Print(y,term)
			req := service.GetCourseRequest{
				UserID:   os.Getenv("USER_ID"),
				Password: os.Getenv("PASSWORD"),
				Year:     y,
				Term:     term,
			}

			courses, err := service.FetchCourses(req)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("fetched courses: %d\n", len(courses.Classes))

			// save
			if err := repository.SaveCourse(database, courses); err != nil {
				log.Fatal(err)
			}
			if err := repository.SaveClasses(database, courses.Classes); err != nil {
				log.Fatal(err)
			}
		}
	}
	// fetch

	log.Println("saved to DB")

	log.Println("done")
}
