package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
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
