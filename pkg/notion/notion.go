package notion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrMissingNotionAPIKey    = errors.New("NOTION_API_KEY is missing")
	ErrMissingClassDatabaseID = errors.New("NOTION_CLASS_DATABASE_ID is missing")
	ErrMissingEventDatabaseID = errors.New("NOTION_EVENT_DATABASE_ID is missing")
)

var defaultHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

const maxNotionRequestAttempts = 4

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
	req.Header.Set("Authorization", "Bearer "+LoadConfigFromEnv().APIToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	// 実行
	res, body, err := doNotionRequest(req)
	if err != nil {
		return "", fmt.Errorf("request send error: %w", err)
	}
	defer res.Body.Close()

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

	req.Header.Set("Authorization", "Bearer "+LoadConfigFromEnv().APIToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, resBody, err := doNotionRequest(req)
	if err != nil {
		return fmt.Errorf("request send error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("notion update error: status=%d body=%s", res.StatusCode, string(resBody))
	}

	return nil
}

func QueryDatabase(databaseID string, payload any) (*databaseQueryResponse, error) {
	url := "https://api.notion.com/v1/databases/" + databaseID + "/query"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request create error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+LoadConfigFromEnv().APIToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	res, body, err := doNotionRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request send error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("notion query error: status=%d body=%s", res.StatusCode, string(body))
	}

	var result databaseQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &result, nil
}

func doNotionRequest(req *http.Request) (*http.Response, []byte, error) {
	var lastErr error

	for attempt := 1; attempt <= maxNotionRequestAttempts; attempt++ {
		clonedReq := cloneRequest(req)

		res, err := defaultHTTPClient.Do(clonedReq)
		if err != nil {
			lastErr = err
			if attempt == maxNotionRequestAttempts {
				return nil, nil, err
			}
			time.Sleep(backoffDuration(attempt))
			continue
		}

		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			lastErr = readErr
			if attempt == maxNotionRequestAttempts {
				return nil, nil, readErr
			}
			time.Sleep(backoffDuration(attempt))
			continue
		}

		if shouldRetryStatus(res.StatusCode) && attempt < maxNotionRequestAttempts {
			time.Sleep(retryDelay(res, attempt))
			lastErr = fmt.Errorf("notion retryable status=%d body=%s", res.StatusCode, string(body))
			continue
		}

		res.Body = io.NopCloser(bytes.NewReader(body))
		return res, body, nil
	}

	return nil, nil, lastErr
}

func cloneRequest(req *http.Request) *http.Request {
	cloned := req.Clone(req.Context())
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err == nil {
			cloned.Body = body
		}
	}
	return cloned
}

func shouldRetryStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= 500
}

func retryDelay(res *http.Response, attempt int) time.Duration {
	if header := res.Header.Get("Retry-After"); header != "" {
		if seconds, err := strconv.Atoi(header); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return backoffDuration(attempt)
}

func backoffDuration(attempt int) time.Duration {
	return time.Duration(attempt) * 500 * time.Millisecond
}
