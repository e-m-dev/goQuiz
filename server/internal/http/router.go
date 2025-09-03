package http

import (
	"log"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	origins := os.Getenv("CORS_ORIGINS")
	allowed := strings.Split(origins, ",")

	if len(allowed) == 1 && allowed[0] == "" {
		allowed = nil
		log.Printf("Incorrect ENV file setup -> no AllowedOrigins defined")
	}

	log.Printf("CORS AllowedOrigins=%v", allowed)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowed,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.With(rateLimitCreate).Post("/rooms", h.CreateRoomHandler)
	r.Get("/rooms/{code}", h.GetRoomHandler)
	r.With(rateLimitJoin).Post("/rooms/{code}/join", h.JoinRoomHandler)
	r.Post("/rooms/{code}/leave", h.LeaveRoomHandler)
	r.Get("/ws/{code}", h.WSHandler)
	r.Post("/questions", h.CreateQuestionHandler)
	r.Post("/questions/gen", h.GenerateQuestionsHandler)
	r.Get("/questions/rnd", h.FetchRandomQuestionsHandler)

	return r
}
