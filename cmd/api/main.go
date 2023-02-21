// package main

// import (
// 	"go-chi-api/api/post"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// )

// func main() {
// 	port := "8080" // TODO read from env
// 	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
// 	log.Printf("POSTGRES_HOST: %s", POSTGRES_HOST)
// 	// database.Connect()

// 	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
// 		port = fromEnv
// 	}

// 	log.Printf(">>>Stesttarting server on http://localhost:%s", port)

// 	r := chi.NewRouter()
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	// r.Use(middleware.Timeout(60 * time.Second))

// 	// health check
// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/plain")
// 		w.Write([]byte("Hello World!"))
// 	})

// 	r.Mount("/posts", post.PostsResource{}.Routes())

// 	log.Fatal(http.ListenAndServe(":"+port, r))
// }
