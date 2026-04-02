package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/db"
	"github.com/med-000/notifyclass/pkg/notion"
)

func main() {
	_ = godotenv.Load()

	//DB接続
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	//env取得
	classDBID := os.Getenv("NOTION_CLASS_DATABASE_ID")
	eventDBID := os.Getenv("NOTION_EVENT_DATABASE_ID")

	if classDBID == "" || eventDBID == "" {
		log.Fatal("NOTION DATABASE ID is missing")
	}

	//class取得 & sync
	var classes []db.Class

	if err := dbConn.Find(&classes).Error; err != nil {
		log.Fatal("class fetch error:", err)
	}

	log.Println("start class sync...")

	if err := notion.SyncClassesToNotion(dbConn, classes, classDBID); err != nil {
		log.Fatal("class sync error:", err)
	}

	log.Println("class sync done")

	//event JOIN取得
	var events []notion.SyncEvent

	err = dbConn.
		Table("events").
		Select(`
			events.external_id as event_external_id,
			events.name,
			events.group,
			events.category,
			events.start_at,
			classes.external_id as class_external_id
		`).
		Joins("JOIN classes ON events.class_id = classes.id").
		Scan(&events).Error

	if err != nil {
		log.Fatal("event join error:", err)
	}

	log.Println("start event sync...")

	// event sync
	if err := notion.SyncEventsToNotion(dbConn, events, eventDBID); err != nil {
		log.Fatal("event sync error:", err)
	}

	log.Println("event sync done")
	log.Println("all done")
}