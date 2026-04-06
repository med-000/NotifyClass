package main

import (
	"fmt"
	"log"
	"os"
	"time"

	discordgo "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/appflow"
	appdiscord "github.com/med-000/notifyclass/pkg/discord"
	appLogger "github.com/med-000/notifyclass/pkg/logger"
)

type botHandler struct {
	log     *appLogger.DiscordLogger
	service *appdiscord.Service
}

func main() {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	discordLogger, err := appLogger.NewDiscordLogger()
	if err != nil {
		log.Fatal(err)
	}

	repositoryLogger, err := appLogger.NewRepositoryLogger()
	if err != nil {
		log.Fatal(err)
	}

	service := appdiscord.NewService(dbConn, discordLogger, repositoryLogger)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	handler := &botHandler{
		log:     discordLogger,
		service: service,
	}

	dg.AddHandler(handler.onReady)
	dg.AddHandler(handler.onMessage)

	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	discordLogger.Info.Printf("discord bot is running")
	select {}
}

func (h *botHandler) onReady(_ *discordgo.Session, r *discordgo.Ready) {
	h.log.Info.Printf("ready user=%s", r.User.String())
}

func (h *botHandler) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	if s.State == nil || s.State.User == nil {
		h.log.Warn.Printf("discord state user is nil")
		return
	}

	content := m.Content

	switch {
	case content == "/sync":
		h.log.Info.Printf("sync command received author=%s channel=%s", m.Author.Username, m.ChannelID)
		h.replyWithSyncRequest(s, m.ChannelID, m.Author.Username)
	case isMentioned(m.Mentions, s.State.User.ID):
		h.log.Info.Printf("mention received author=%s channel=%s", m.Author.Username, m.ChannelID)
		if h.replyIfBusy(s, m.ChannelID) {
			return
		}
		h.replyWithPendingTasks(s, m.ChannelID)
	case appdiscord.IsSlotCommand(content):
		h.log.Info.Printf("slot command received author=%s channel=%s content=%s", m.Author.Username, m.ChannelID, content)
		if h.replyIfBusy(s, m.ChannelID) {
			return
		}
		h.replyWithSlotTasks(s, m.ChannelID, content)
	}
}

func (h *botHandler) replyWithSyncRequest(s *discordgo.Session, channelID string, username string) {
	if h.replyIfBusy(s, channelID) {
		return
	}

	accepted, err := appflow.RequestSync("discord:" + username)
	if err != nil {
		h.log.Error.Printf("failed to request sync err=%v", err)
		if _, sendErr := s.ChannelMessageSend(channelID, "sync リクエストの受付に失敗しました。"); sendErr != nil {
			h.log.Error.Printf("send error channel=%s err=%v", channelID, sendErr)
		}
		return
	}

	message := "sync リクエストを受け付けました。backend で順次処理します。"
	if !accepted {
		message = "現在 sync は実行中、またはキュー済みです。処理完了まで待ってください。"
	}

	if _, err := s.ChannelMessageSend(channelID, message); err != nil {
		h.log.Error.Printf("send error channel=%s err=%v", channelID, err)
		return
	}

	h.log.Info.Printf("sync reply sent channel=%s accepted=%t", channelID, accepted)
}

func (h *botHandler) replyIfBusy(s *discordgo.Session, channelID string) bool {
	status, err := appflow.ReadStatus()
	if err != nil {
		h.log.Error.Printf("failed to read appflow status err=%v", err)
		return false
	}

	if status == nil || !status.Busy {
		return false
	}

	message := fmt.Sprintf(
		"現在 `%s` を実行中です。処理が終わってからもう一度試してください。開始時刻: %s",
		status.Stage,
		status.StartedAt.Format("2006-01-02 15:04:05"),
	)

	if _, err := s.ChannelMessageSend(channelID, message); err != nil {
		h.log.Error.Printf("send busy status error channel=%s err=%v", channelID, err)
	}

	return true
}

func (h *botHandler) replyWithPendingTasks(s *discordgo.Session, channelID string) {
	message, err := h.service.BuildPendingTasksMessage(time.Now())
	if err != nil {
		h.log.Error.Printf("failed to build pending tasks message err=%v", err)
		message = "未完了タスク一覧を取得できませんでした。"
	}

	if _, err := s.ChannelMessageSend(channelID, message); err != nil {
		h.log.Error.Printf("send error channel=%s err=%v", channelID, err)
		return
	}

	h.log.Info.Printf("reply sent channel=%s", channelID)
}

func (h *botHandler) replyWithSlotTasks(s *discordgo.Session, channelID string, content string) {
	message, err := h.service.BuildPendingTasksBySlotMessage(content)
	if err != nil {
		h.log.Warn.Printf("failed to build slot tasks message content=%s err=%v", content, err)
		message = "コマンド形式が不正です。`/2025010102` のように `/YYYYTTDDPP` で送信してください。"
	}

	if _, err := s.ChannelMessageSend(channelID, message); err != nil {
		h.log.Error.Printf("send error channel=%s err=%v", channelID, err)
		return
	}

	h.log.Info.Printf("reply sent channel=%s", channelID)
}

func isMentioned(users []*discordgo.User, botUserID string) bool {
	for _, user := range users {
		if user.ID == botUserID {
			return true
		}
	}
	return false
}
