package appflow

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

func BuildFetchRequests(now time.Time) ([]service.GetCourseRequest, error) {
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	if strings.EqualFold(os.Getenv("APP_INITIAL_SYNC"), "true") {
		done, err := IsInitialSyncAlreadyDone()
		if err != nil {
			return nil, err
		}
		if done {
			year, term, err := resolveTargetAcademicPeriod(now)
			if err != nil {
				return nil, err
			}

			return []service.GetCourseRequest{
				{
					UserID:   userID,
					Password: password,
					Year:     year,
					Term:     term,
				},
			}, nil
		}

		return buildInitialFetchRequests(userID, password, now)
	}

	year, term, err := resolveTargetAcademicPeriod(now)
	if err != nil {
		return nil, err
	}

	return []service.GetCourseRequest{
		{
			UserID:   userID,
			Password: password,
			Year:     year,
			Term:     term,
		},
	}, nil
}

func FetchCourses(now time.Time) ([]*parser.Course, error) {
	requests, err := BuildFetchRequests(now)
	if err != nil {
		return nil, err
	}

	s, err := NewService()
	if err != nil {
		return nil, err
	}

	courses := make([]*parser.Course, 0, len(requests))
	for _, req := range requests {
		course, err := s.FetchAll(req)
		if err != nil {
			return nil, fmt.Errorf("fetch year=%d term=%d: %w", req.Year, req.Term, err)
		}
		courses = append(courses, course)
	}

	return courses, nil
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

func SaveCoursesToDB(dbConn *gorm.DB, courses []*parser.Course) error {
	s, err := NewService()
	if err != nil {
		return err
	}

	for _, course := range courses {
		if err := s.SaveAll(dbConn, course); err != nil {
			return err
		}
	}

	return nil
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

func resolveTargetAcademicPeriod(now time.Time) (int, int, error) {
	year := inferAcademicYear(now)
	term := inferAcademicTerm(now)

	if value := os.Getenv("APP_YEAR"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid APP_YEAR: %w", err)
		}
		year = parsed
	}

	if value := os.Getenv("APP_TERM"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid APP_TERM: %w", err)
		}
		if parsed != 1 && parsed != 2 {
			return 0, 0, fmt.Errorf("APP_TERM must be 1 or 2")
		}
		term = parsed
	}

	return year, term, nil
}

func buildInitialFetchRequests(userID, password string, now time.Time) ([]service.GetCourseRequest, error) {
	startYear, err := enrollmentYearFromUserID(userID)
	if err != nil {
		return nil, err
	}

	endYear, endTerm, err := resolveTargetAcademicPeriod(now)
	if err != nil {
		return nil, err
	}

	requests := make([]service.GetCourseRequest, 0, (endYear-startYear+1)*2)
	for year := startYear; year <= endYear; year++ {
		for term := 1; term <= 2; term++ {
			if year == endYear && term > endTerm {
				break
			}

			requests = append(requests, service.GetCourseRequest{
				UserID:   userID,
				Password: password,
				Year:     year,
				Term:     term,
			})
		}
	}

	return requests, nil
}

func enrollmentYearFromUserID(userID string) (int, error) {
	if len(userID) < 2 {
		return 0, fmt.Errorf("USER_ID is too short to infer enrollment year")
	}

	yearPrefix := userID[:2]
	year, err := strconv.Atoi(yearPrefix)
	if err != nil {
		return 0, fmt.Errorf("invalid USER_ID enrollment year: %w", err)
	}

	return 2000 + year, nil
}

func inferAcademicTerm(now time.Time) int {
	month := now.Month()
	if month >= 4 && month <= 9 {
		return 1
	}
	return 2
}

func inferAcademicYear(now time.Time) int {
	if now.Month() >= 1 && now.Month() <= 3 {
		return now.Year() - 1
	}
	return now.Year()
}
