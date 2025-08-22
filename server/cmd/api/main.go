package main

import (
	internHttp "goQuiz/server/internal/http"
	"goQuiz/server/internal/store"
	wshub "goQuiz/server/internal/ws"
	"log"
	"net/http"
	"os"

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

	log.Printf("server running on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
