package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/notion"
)

func main() {
	_ = godotenv.Load()

	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	notionLogger, err := logger.NewNotionLogger()
	if err != nil {
		log.Fatal(err)
	}

	cfg := notion.LoadConfigFromEnv()

	if err := notion.SyncEventCompletionFromNotion(dbConn, notionLogger, cfg); err != nil {
		log.Fatal("notion pull error:", err)
	}
}
