package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	godotenv.Load(".env")
	r := chi.NewRouter()

	bind := os.Getenv("QM_BIND")
	port := os.Getenv("QM_PORT")
	addr := bind + ":" + port

	token := os.Getenv("QM_API")
	if token == "" {
		log.Fatal("API Token Missing")
	}

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.With(bearerAuth(token)).Post("/v1/questions/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gen"))
	})

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
