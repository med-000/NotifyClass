package appflow

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/notion"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/service"
	"gorm.io/gorm"
)

func NewService() (*service.Service, error) {
	serviceLogger, err := logger.NewServiceLogger()
	if err != nil {
		return nil, err
	}

	return service.NewService(serviceLogger), nil
}

func NewNotionLogger() (*logger.NotionLogger, error) {
	return logger.NewNotionLogger()
}

func DefaultFetchRequest() service.GetCourseRequest {
	return service.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     2,
	}
}

func FetchCourse() (*parser.Course, error) {
	s, err := NewService()
	if err != nil {
		return nil, err
	}

	return s.FetchAll(DefaultFetchRequest())
}

func ExportCourseToJSON(filename string, course any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(course)
}

func SaveCourseToDB(dbConn *gorm.DB, course *parser.Course) error {
	s, err := NewService()
	if err != nil {
		return err
	}

	return s.SaveAll(dbConn, course)
}

func SyncNotionPull(dbConn *gorm.DB) error {
	notionLogger, err := NewNotionLogger()
	if err != nil {
		return err
	}

	cfg := notion.LoadConfigFromEnv()
	notionLogger.Info.Printf("start notion pull")
	if err := notion.SyncEventCompletionFromNotion(dbConn, notionLogger, cfg); err != nil {
		return fmt.Errorf("sync notion pull: %w", err)
	}
	notionLogger.Info.Printf("notion pull done")

	return nil
}

func SyncNotionPush(dbConn *gorm.DB) error {
	notionLogger, err := NewNotionLogger()
	if err != nil {
		return err
	}

	cfg := notion.LoadConfigFromEnv()
	notionLogger.Info.Printf("start notion push")
	if err := notion.SyncAllFromDB(dbConn, notionLogger, cfg); err != nil {
		return fmt.Errorf("sync notion push: %w", err)
	}
	notionLogger.Info.Printf("notion push done")

	return nil
}
