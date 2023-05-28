package api

import (
	"go-chi-api/internal/api/v1/customer"
	"go-chi-api/internal/api/v1/item"
	"go-chi-api/internal/api/v1/item_type"
	"go-chi-api/internal/api/v1/login"
	"go-chi-api/internal/api/v1/order"
	"go-chi-api/internal/api/v1/restock"
	"go-chi-api/internal/api/v1/user"
	"go-chi-api/internal/database"
	"go-chi-api/internal/migrations"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(identifier string, port string) {
	migrations.Migrate()
	database.Connect()
	log.Printf("Starting server on http://localhost:%s", port)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello World!"))
	})

	r.Mount("/item", item.ItemsResource{}.Routes())

	r.Mount("/login", login.LoginResource{}.Routes())

	r.Mount("/user", user.UsersResource{}.Routes())

	r.Mount("/item_type", item_type.ItemTypesResource{}.Routes())

	r.Mount("/restock", restock.RestockResource{}.Routes())

	r.Mount("/customer", customer.CustomerResource{}.Routes())

	r.Mount("/order", order.OrderResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
