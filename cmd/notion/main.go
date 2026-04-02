package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/notion"
	"github.com/med-000/notifyclass/pkg/repository"
)

func main() {
	_ = godotenv.Load()

	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	databaseID := os.Getenv("NOTION_DATABASE_ID")
	if databaseID == "" {
		log.Fatal("NOTION_DATABASE_ID is empty")
	}

	var classes []db.Class
if err := dbConn.Find(&classes).Error; err != nil {
	log.Fatal("db fetch error:", err)
}

for _, c := range classes {
	payload := notion.BuildClassPayload(c)
	if payload == nil {
		continue
	}

	pageID, err := repository.UpsertNotion(
		dbConn,
		c.ExternalID,
		db.MappingTypeClass,
		payload,
	)
	if err != nil {
		log.Println("sync error:", err)
		continue
	}

	log.Println("synced:", c.Title, pageID)
}
}
