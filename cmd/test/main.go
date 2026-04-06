package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/notion"
	"github.com/med-000/notifyclass/pkg/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("なんかおかしい")
	}
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	log.Println("DB ready")

	serviceLogger, _ := logger.NewServiceLogger()
	s := service.NewService(serviceLogger)
	notionLogger, _ := logger.NewNotionLogger()

	req := service.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     2,
	}
	course, err := s.FetchAll(req)
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
	err = s.SaveAll(database, course)
	if err != nil {
		fmt.Println("failed to save database:", err)
		return
	}

	cfg := notion.LoadConfigFromEnv()
	err = notion.SyncAllFromDB(database, notionLogger, cfg)
	if err != nil {
		fmt.Println("failed to sync notion:", err)
		return
	}

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
