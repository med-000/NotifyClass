package appflow

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

const (
	statusFilePath      = "runtime/appflow_status.json"
	syncRequestFilePath = "runtime/sync_request.json"
	defaultSyncSchedule = "0 6 * * *"
)

type Status struct {
	Busy      bool      `json:"busy"`
	Stage     string    `json:"stage"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastError string    `json:"last_error,omitempty"`
}

type SyncRequest struct {
	Source      string    `json:"source"`
	RequestedAt time.Time `json:"requested_at"`
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

	courses, err := fetchCoursesWithStatus()
	if err != nil {
		return failPipeline("fetch", err)
	}

	if err := updateStage("json export"); err != nil {
		return err
	}
	if err := ExportCourseToJSON(exportFilename, courses); err != nil {
		return failPipeline("json export", err)
	}

	if err := updateStage("db save"); err != nil {
		return err
	}
	if err := SaveCoursesToDB(dbConn, courses); err != nil {
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

func RunScheduler(dbConn *gorm.DB, exportFilename string) error {
	schedule := os.Getenv("APP_SYNC_SCHEDULE")
	if schedule == "" {
		schedule = defaultSyncSchedule
	}

	triggerCh := make(chan string, 1)
	scheduler := cron.New()

	if _, err := scheduler.AddFunc(schedule, func() {
		enqueueTrigger(triggerCh, "scheduled")
	}); err != nil {
		return fmt.Errorf("add cron schedule: %w", err)
	}

	scheduler.Start()
	defer scheduler.Stop()

	log.Printf("backend scheduler started schedule=%s", schedule)

	requestTicker := time.NewTicker(3 * time.Second)
	defer requestTicker.Stop()

	for {
		select {
		case source := <-triggerCh:
			if err := runTriggeredPipeline(dbConn, exportFilename, source); err != nil {
				log.Printf("pipeline run failed source=%s err=%v", source, err)
			}
		case <-requestTicker.C:
			request, err := ConsumeSyncRequest()
			if err != nil {
				log.Printf("failed to consume sync request err=%v", err)
				continue
			}
			if request == nil {
				continue
			}
			enqueueTrigger(triggerCh, request.Source)
		}
	}
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

func ReadSyncRequest() (*SyncRequest, error) {
	body, err := os.ReadFile(syncRequestFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var request SyncRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

func RequestSync(source string) (bool, error) {
	status, err := ReadStatus()
	if err != nil {
		return false, err
	}
	if status != nil && status.Busy {
		return false, nil
	}

	request, err := ReadSyncRequest()
	if err != nil {
		return false, err
	}
	if request != nil {
		return false, nil
	}

	request = &SyncRequest{
		Source:      source,
		RequestedAt: time.Now(),
	}

	if err := os.MkdirAll(filepath.Dir(syncRequestFilePath), 0o755); err != nil {
		return false, err
	}

	body, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return false, err
	}

	if err := os.WriteFile(syncRequestFilePath, body, 0o644); err != nil {
		return false, err
	}

	return true, nil
}

func ConsumeSyncRequest() (*SyncRequest, error) {
	request, err := ReadSyncRequest()
	if err != nil || request == nil {
		return request, err
	}

	if err := os.Remove(syncRequestFilePath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return request, nil
}

func fetchCoursesWithStatus() ([]*parser.Course, error) {
	if err := updateStage("fetch"); err != nil {
		return nil, err
	}

	return FetchCourses(time.Now())
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
	err := os.Remove(statusFilePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
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

func runTriggeredPipeline(dbConn *gorm.DB, exportFilename, source string) error {
	status, err := ReadStatus()
	if err != nil {
		return err
	}
	if status != nil && status.Busy {
		log.Printf("skip trigger source=%s reason=busy stage=%s", source, status.Stage)
		return nil
	}

	log.Printf("start pipeline source=%s", source)
	return RunFullPipeline(dbConn, exportFilename)
}

func enqueueTrigger(triggerCh chan<- string, source string) {
	select {
	case triggerCh <- source:
	default:
		log.Printf("skip enqueue trigger source=%s reason=already queued", source)
	}
}
