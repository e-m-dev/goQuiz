package main

import (
	"goQuiz/server/internal/cfg"
	internHttp "goQuiz/server/internal/http"
	"goQuiz/server/internal/store"
	wshub "goQuiz/server/internal/ws"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	s := store.NewStore()
	hub := wshub.NewHub()
	h := &internHttp.Handler{Ref: s, Hub: hub}

	r := internHttp.NewRouter(h)

	bind := getEnv("BIND_ADDR", "0.0.0.0")
	port := getEnv("PORT", "8080")
	addr := bind + ":" + port

	cfg.Debug = strings.EqualFold(os.Getenv("DEBUG"), "true")

	cfg.DBPath = getEnv("DB_PATH", "./data/q.db")
	log.Printf("DB Path Found -> %s", cfg.DBPath)

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

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
