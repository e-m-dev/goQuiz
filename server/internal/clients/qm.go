package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goQuiz/server/internal/store"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL string, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Fetch(ctx context.Context, count int, category *string, difficulty *string) ([]store.Question, error) {
	url := strings.TrimRight(c.baseURL, "/") + "/v1/questions/generate"

	body := map[string]any{"count": count}
	if category == nil {
		body["category"] = strings.TrimSpace(*category)
	}
	if difficulty == nil {
		body["difficulty"] = strings.TrimSpace(*difficulty)
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qm %d %s", resp.StatusCode, string(b))
	}

	var payload struct {
		Questions []struct {
			Prompt       string   `json:"prompt"`
			Options      []string `json:"options"`
			CorrectIndex int      `json:"correctIndex"`
			Category     *string  `json:"category"`
			Difficulty   *string  `json:"difficulty"`
		} `json:"questions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	res := make([]store.Question, 0, len(payload.Questions))
	for _, q := range payload.Questions {
		res = append(res, store.Question{
			Prompt:       strings.TrimSpace(q.Prompt),
			Options:      q.Options,
			CorrectIndex: q.CorrectIndex,
			Category:     q.Category,
			Difficulty:   q.Difficulty,
		})
	}

	return res, nil
}
