package discord

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/repository"
	"gorm.io/gorm"
)

const maxDisplayedTasks = 10

type Service struct {
	dbConn        *gorm.DB
	log           *logger.DiscordLogger
	repositoryLog *logger.RepositoryLogger
}

func NewService(dbConn *gorm.DB, discordLog *logger.DiscordLogger, repositoryLog *logger.RepositoryLogger) *Service {
	return &Service{
		dbConn:        dbConn,
		log:           discordLog,
		repositoryLog: repositoryLog,
	}
}

func (s *Service) BuildPendingTasksMessage(now time.Time) (string, error) {
	eventRepo := repository.NewEventRepository(s.dbConn, s.repositoryLog)
	tasks, err := eventRepo.FindPendingTasks()
	if err != nil {
		s.log.Error.Printf("failed FindPendingTasks err=%v", err)
		return "", err
	}

	s.log.Info.Printf("build pending tasks message count=%d", len(tasks))
	return formatPendingTasksMessage(
		fmt.Sprintf("未完了タスク一覧です。(%s)", now.Format("2006-01-02")),
		selectTasksForDisplay(tasks, now),
		len(tasks),
	), nil
}

func (s *Service) BuildPendingTasksBySlotMessage(content string) (string, error) {
	slot, err := ParseSlotCommand(content)
	if err != nil {
		s.log.Warn.Printf("invalid slot command content=%s err=%v", content, err)
		return "", err
	}

	eventRepo := repository.NewEventRepository(s.dbConn, s.repositoryLog)
	tasks, err := eventRepo.FindPendingTasksBySlot(slot.Year, slot.Term, slot.Day, slot.Period)
	if err != nil {
		s.log.Error.Printf(
			"failed FindPendingTasksBySlot year=%d term=%d day=%d period=%d err=%v",
			slot.Year, slot.Term, slot.Day, slot.Period, err,
		)
		return "", err
	}

	s.log.Info.Printf(
		"build pending tasks by slot year=%d term=%d day=%d period=%d count=%d",
		slot.Year, slot.Term, slot.Day, slot.Period, len(tasks),
	)

	title := fmt.Sprintf(
		"%d年度 %s %s %d限 の未完了タスクです。",
		slot.Year,
		termToString(slot.Term),
		dayToString(slot.Day),
		slot.Period,
	)

	return formatPendingTasksMessage(title, selectTasksForDisplay(tasks, time.Now()), len(tasks)), nil
}

func formatPendingTasksMessage(title string, tasks []repository.PendingTaskRow, totalCount int) string {
	if len(tasks) == 0 {
		return title + "\n該当する未完了タスクはありません。"
	}

	lines := []string{title}
	for i, task := range tasks {
		lines = append(lines, formatTaskLine(i+1, task))
	}
	if totalCount > len(tasks) {
		lines = append(lines, fmt.Sprintf("...他 %d 件", totalCount-len(tasks)))
	}
	return strings.Join(lines, "\n")
}

func selectTasksForDisplay(tasks []repository.PendingTaskRow, now time.Time) []repository.PendingTaskRow {
	upcoming := make([]repository.PendingTaskRow, 0, len(tasks))
	overdue := make([]repository.PendingTaskRow, 0, len(tasks))
	undated := make([]repository.PendingTaskRow, 0, len(tasks))

	for _, task := range tasks {
		if task.Deadline == nil {
			undated = append(undated, task)
			continue
		}

		if task.Deadline.Before(now) {
			overdue = append(overdue, task)
			continue
		}

		upcoming = append(upcoming, task)
	}

	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].Deadline.Before(*upcoming[j].Deadline)
	})

	sort.Slice(overdue, func(i, j int) bool {
		return overdue[i].Deadline.After(*overdue[j].Deadline)
	})

	sort.Slice(undated, func(i, j int) bool {
		if undated[i].Year != undated[j].Year {
			return undated[i].Year < undated[j].Year
		}
		if undated[i].Term != undated[j].Term {
			return undated[i].Term < undated[j].Term
		}
		if undated[i].Day != undated[j].Day {
			return undated[i].Day < undated[j].Day
		}
		return undated[i].Period < undated[j].Period
	})

	switch {
	case len(upcoming) > 0:
		return limitTasks(upcoming, maxDisplayedTasks)
	case len(overdue) > 0:
		return limitTasks(overdue, maxDisplayedTasks)
	default:
		return limitTasks(undated, maxDisplayedTasks)
	}
}

func limitTasks(tasks []repository.PendingTaskRow, max int) []repository.PendingTaskRow {
	if len(tasks) <= max {
		return tasks
	}
	return tasks[:max]
}

func formatTaskLine(index int, task repository.PendingTaskRow) string {
	deadline := "期限なし"
	if task.Deadline != nil {
		deadline = task.Deadline.Format("2006-01-02 15:04")
	}

	return fmt.Sprintf(
		"%d. [%s] %s / %s / %s / %s %d限",
		index,
		emptyToFallback(task.Category, "未分類"),
		task.ClassTitle,
		task.Name,
		deadline,
		dayToString(task.Day),
		task.Period,
	)
}

func dayToString(day int) string {
	switch day {
	case 1:
		return "月曜"
	case 2:
		return "火曜"
	case 3:
		return "水曜"
	case 4:
		return "木曜"
	case 5:
		return "金曜"
	case 6:
		return "土曜"
	case 7:
		return "日曜"
	default:
		return fmt.Sprintf("曜日%d", day)
	}
}

func termToString(term int) string {
	switch term {
	case 1:
		return "前期"
	case 2:
		return "後期"
	default:
		return fmt.Sprintf("学期%d", term)
	}
}

func emptyToFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
