package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/appflow"
)

func main() {
	_ = godotenv.Load()

	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := appflow.SyncNotionPush(dbConn); err != nil {
		log.Fatal(err)
	}
}
