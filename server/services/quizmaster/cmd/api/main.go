package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type questionReq struct {
	Count      int
	Category   *string
	Difficulty *string
}

type genReq struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerateConfig struct {
		ResponseMIMEType string `json:"response_mime_type"`
	} `json:"generationConfig"`
}

func main() {
	godotenv.Load(".env")
	r := chi.NewRouter()

	bind := os.Getenv("QM_BIND")
	port := os.Getenv("QM_PORT")
	addr := bind + ":" + port

	token := os.Getenv("QM_TOKEN")
	if token == "" {
		log.Fatal("API Token Missing")
	}

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.With(bearerAuth(token)).Post("/v1/questions/generate", GenerateQHandler)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server running on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func GenerateQHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.Body = http.MaxBytesReader(w, r.Body, 8192)

	var maxErr *http.MaxBytesError
	var req questionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.As(err, &maxErr) {
			http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "invalid req", http.StatusBadRequest)
		return
	}

	req.Count = clampCount(req.Count)
	req.Category = cleanPtr(req.Category)
	req.Difficulty = cleanPtr(req.Difficulty)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":         true,
		"count":      req.Count,
		"category":   req.Category,
		"difficulty": req.Difficulty,
	})

	log.Printf("QuizMaster Generate echo count {%d} questions in {%s}", req.Count, time.Since(start))

}

func callQuizMaster(ctx context.Context, model string, key string, count int, cat *string, diff *string) (jsonString string, err error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)

	prompt := fmt.Sprintf(`Output ONLY JSON:
		{ "questions":[{ "prompt": string, "options": [string], "correctIndex": number,
		"category": %s, "difficulty": %s, "source": "quizmaster:gemini" }] }

		Rules: exactly %d questions; 4 options per question; short prompts; no trick Qs; correctIndex index of answer.
		No extra text, no markdown. Strictly JSON only.`, valOrNull(cat), valOrNull(diff), count)

	reqBody := genReq{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{Parts: []struct {
				Text string `json:"text"`
			}{{Text: prompt}}},
		},
	}
	reqBody.GenerateConfig.ResponseMIMEType = "application/json"

	b, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Failed to marshal req body")
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		log.Printf("API communication construct failed")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send req")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gem error, code: %d, resp: %s", resp.StatusCode, string(b))
	}

	var gr struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		log.Printf("failed to decode resp body into structure")
		return "", err
	}

	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return gr.Candidates[0].Content.Parts[0].Text, nil
}

func clampCount(original int) int {
	if original <= 50 && original >= 1 {
		return original
	}

	if original > 50 {
		return 50
	}

	return 1

}

func cleanPtr(p *string) *string {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)

	if s == "" {
		return nil
	}

	return &s
}

func valOrNull(s *string) string {
	if s == nil || strings.TrimSpace(*s) == "" {
		return "null"
	}

	return strconv.Quote(strings.TrimSpace(*s))
}

func bearerAuth(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") ||
				strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")) != expected {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
