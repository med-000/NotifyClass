package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	dbConn, err := db.NewDB()
	if err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("db connect error: %v", err))
		log.Fatal(err)
	}

	if err := appflow.SyncNotionPush(dbConn); err != nil {
		_ = appflow.NotifyDiscordError(fmt.Sprintf("notion push error: %v", err))
		log.Fatal(err)
	}
}
