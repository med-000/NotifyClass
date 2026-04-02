package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateNotion(payload any) (string, error) {
	url := "https://api.notion.com/v1/pages"

	// JSON化
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request create error: %w", err)
	}

	// ヘッダー
	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	// 実行
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request send error: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	// エラーハンドリング
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("notion api error: status=%d body=%s", res.StatusCode, string(body))
	}

	// レスポンスパース
	var result struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("json unmarshal error: %w", err)
	}

	return result.ID, nil
}

func UpdateNotion(pageID string, payload any) error {
	url := "https://api.notion.com/v1/pages/" + pageID

	bodyMap, ok := payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format")
	}

	properties, ok := bodyMap["properties"]
	if !ok {
		return fmt.Errorf("properties not found in payload")
	}

	reqBody := map[string]interface{}{
		"properties": properties,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request create error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request send error: %w", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 300 {
		return fmt.Errorf("notion update error: status=%d body=%s", res.StatusCode, string(resBody))
	}

	return nil
}
