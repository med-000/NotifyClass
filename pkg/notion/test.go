package notion

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Notion() {
	_ = godotenv.Load()

	databaseID := os.Getenv("NOTION_DATABASE_ID")
	if databaseID == "" {
		panic("NOTION_DATABASE_ID is empty")
	}

	url := "https://api.notion.com/v1/databases/" + databaseID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("status:", res.Status)

	// ===== ファイル出力 =====
	err = os.WriteFile("notion_debug.json", body, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("saved to notion_debug.json")
}

func CreatePage() {
	dataBaseID := os.Getenv("NOTION_DATABASE_ID")
	token := os.Getenv("NOTION_API_KEY")

	url := "https://api.notion.com/v1/pages"

	jsonBody := `{
		"parent": { "database_id": "` + dataBaseID + `" },
		"properties": {
			"Name": {
				"title": [
					{
						"text": { "content": "テストイベント" }
					}
				]
			}
		}
	}`

	req, err := http.NewRequest("POST", url, strings.NewReader(jsonBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println("status:", res.Status)
	fmt.Println(string(body))
}
