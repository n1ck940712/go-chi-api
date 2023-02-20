package api

import (
	"go-chi-api/api/post"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(identifier string, port string) {

	log.Printf("Starting server on http://localhost:%s", port)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// r.Use(middleware.Timeout(60 * time.Second))

	// health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello World!"))
	})

	r.Mount("/posts", post.PostsResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}