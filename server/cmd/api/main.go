package main

import (
	internHttp "goQuiz/server/internal/http"
	"goQuiz/server/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	s := store.NewStore()
	h := &internHttp.Handler{Ref: s}

	r := internHttp.NewRouter(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on :%s", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, r); err != nil {
	}
}
