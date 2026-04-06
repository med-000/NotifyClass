package main

import (
	"log"
	"os"
	"time"

	discordgo "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
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
	case isMentioned(m.Mentions, s.State.User.ID):
		h.log.Info.Printf("mention received author=%s channel=%s", m.Author.Username, m.ChannelID)
		h.replyWithPendingTasks(s, m.ChannelID)
	case appdiscord.IsSlotCommand(content):
		h.log.Info.Printf("slot command received author=%s channel=%s content=%s", m.Author.Username, m.ChannelID, content)
		h.replyWithSlotTasks(s, m.ChannelID, content)
	}
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
