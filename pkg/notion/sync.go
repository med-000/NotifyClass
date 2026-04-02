package notion

import (
	"log"
	"time"

	"github.com/med-000/notifyclass/db"
	"gorm.io/gorm"
)

func SyncClassesToNotion(
	dbConn *gorm.DB,
	classes []db.Class,
	databaseID string,
) error {

	for _, c := range classes {
		payload := BuildClassPayload(databaseID, c)
		if payload == nil {
			log.Println("payload build failed:", c.ExternalID)
			continue
		}

		mapping, err := FindMapping(
			dbConn,
			c.ExternalID,
			db.MappingTypeClass,
		)

		if err != nil {
			log.Println("mapping fetch error:", err)
			continue
		}

		if mapping == nil {
			pageID, err := CreateNotion(payload)
			if err != nil {
				log.Println("create notion error:", err)
				continue
			}

			if err := UpsertMapping(
				dbConn,
				c.ExternalID,
				db.MappingTypeClass,
				pageID,
			); err != nil {
				log.Println("mapping save error:", err)
			}

			log.Println("created:", c.Title, pageID)

		} else {
			err := UpdateNotion(mapping.NotionPageID, payload)
			if err != nil {
				log.Println("update notion error:", err)
				continue
			}

			log.Println("updated:", c.Title, mapping.NotionPageID)
		}
	}

	return nil
}

func SyncEventsToNotion(dbConn *gorm.DB, events []SyncEvent, databaseID string) error {

	for _, ev := range events {

		//class mapping取得
		classMapping, err := FindMapping(
			dbConn,
			ev.ClassExternalID,
			db.MappingTypeClass,
		)
		if err != nil {
			log.Println("mapping error:", err)
			continue
		}
		if classMapping == nil {
			log.Println("class mapping not found:", ev.ClassExternalID)
			continue
		}

		// =========================
		// ② payload作成
		// =========================
		payload := BuildEventPayload(
			databaseID,
			EventWithRelation{
				ExternalID:  ev.EventExternalID,
				Name:        ev.Name,
				Group:       ev.Group,
				Category:    ev.Category,
				StartAt:     ev.StartAt,
				ClassPageID: classMapping.NotionPageID,
			},
		)

		if payload == nil {
			log.Println("payload build failed:", ev.EventExternalID)
			continue
		}

		// =========================
		// ③ Notion同期（retry付き）
		// =========================
		var pageID string

		for i := 0; i < 3; i++ { // 最大3回リトライ

			pageID, err = UpsertNotion(
				dbConn,
				ev.EventExternalID,
				db.MappingTypeEvent,
				payload,
			)

			if err == nil {
				break
			}

			log.Println("retry:", i+1, "error:", err)
			time.Sleep(1 * time.Second)
		}

		if err != nil {
			log.Println("event sync failed:", ev.Name, err)
			continue
		}

		log.Println("event synced:", ev.Name, pageID)

		// =========================
		// ④ レート制御
		// =========================
		time.Sleep(350 * time.Millisecond)
	}

	return nil
}