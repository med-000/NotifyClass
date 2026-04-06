package appflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/med-000/notifyclass/pkg/parser"
	"gorm.io/gorm"
)

const statusFilePath = "runtime/appflow_status.json"

type Status struct {
	Busy      bool      `json:"busy"`
	Stage     string    `json:"stage"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastError string    `json:"last_error,omitempty"`
}

func RunFullPipeline(dbConn *gorm.DB, exportFilename string) error {
	if err := setStatus(Status{
		Busy:      true,
		Stage:     "notion pull",
		StartedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	cleanup := true
	defer func() {
		if cleanup {
			_ = clearStatus()
		}
	}()

	if err := SyncNotionPull(dbConn); err != nil {
		return failPipeline("notion pull", err)
	}

	course, err := fetchCourseWithStatus()
	if err != nil {
		return failPipeline("fetch", err)
	}

	if err := updateStage("json export"); err != nil {
		return err
	}
	if err := ExportCourseToJSON(exportFilename, course); err != nil {
		return failPipeline("json export", err)
	}

	if err := updateStage("db save"); err != nil {
		return err
	}
	if err := SaveCourseToDB(dbConn, course); err != nil {
		return failPipeline("db save", err)
	}

	if err := updateStage("notion push"); err != nil {
		return err
	}
	if err := SyncNotionPush(dbConn); err != nil {
		return failPipeline("notion push", err)
	}

	cleanup = false
	return clearStatus()
}

func ReadStatus() (*Status, error) {
	body, err := os.ReadFile(statusFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Status{}, nil
		}
		return nil, err
	}

	var status Status
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func fetchCourseWithStatus() (*parser.Course, error) {
	if err := updateStage("fetch"); err != nil {
		return nil, err
	}

	return FetchCourse()
}

func updateStage(stage string) error {
	status, err := ReadStatus()
	if err != nil {
		return err
	}

	if status.StartedAt.IsZero() {
		status.StartedAt = time.Now()
	}

	status.Busy = true
	status.Stage = stage
	status.UpdatedAt = time.Now()
	return setStatus(*status)
}

func failPipeline(stage string, err error) error {
	message := fmt.Sprintf("%s error: %v", stage, err)

	status, readErr := ReadStatus()
	if readErr == nil {
		status.Busy = false
		status.Stage = stage
		status.LastError = message
		status.UpdatedAt = time.Now()
		_ = setStatus(*status)
	}

	_ = NotifyDiscordError(message)
	return err
}

func setStatus(status Status) error {
	if err := os.MkdirAll(filepath.Dir(statusFilePath), 0o755); err != nil {
		return err
	}

	body, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statusFilePath, body, 0o644)
}

func clearStatus() error {
	return os.Remove(statusFilePath)
}

func NotifyDiscordError(message string) error {
	token := os.Getenv("DISCORD_TOKEN")
	channelID := os.Getenv("DISCORD_NOTIFY_CHANNEL_ID")
	if token == "" || channelID == "" {
		return nil
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	_, err = session.ChannelMessageSend(channelID, message)
	return err
}
