package item

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"time"

	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"

	"github.com/go-chi/chi/v5"
)

type ItemsResource struct{}

func (rs ItemsResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.JwtAuthentication)
	r.Get("/", rs.List)
	r.Post("/", rs.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(ctx)
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs ItemsResource) List(w http.ResponseWriter, r *http.Request) {
	result := []models.ItemTable{}
	database.DB.Find(&result)
	response.JSON(w, http.StatusOK, result)
}

func (rs ItemsResource) Create(w http.ResponseWriter, r *http.Request) {
	item := models.ItemTable{}
	json.NewDecoder(r.Body).Decode(&item)
	database.DB.Create(&item)
	response.JSON(w, http.StatusOK, item)
}

func ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs ItemsResource) Get(w http.ResponseWriter, r *http.Request) {
	result := models.ItemTable{}
	id := r.Context().Value("id").(string)
	database.DB.Find(&result, id)
	response.JSON(w, http.StatusOK, result)

}

func (rs ItemsResource) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	item := models.ItemTable{}
	json.NewDecoder(r.Body).Decode(&item)
	database.DB.Model(&item).Where("id = ?", id).UpdateColumns(
		map[string]interface{}{
			"name":       item.Name,
			"attribute":  item.Attribute,
			"updated_at": time.Now(),
		},
	)
	updated_item := models.ItemTable{}
	database.DB.Find(&updated_item, id)
	response.JSON(w, http.StatusOK, updated_item)
}

func (rs ItemsResource) Delete(w http.ResponseWriter, r *http.Request) {
	item := models.ItemTable{}
	id := r.Context().Value("id").(string)
	database.DB.Delete(&item, id)
	response.JSON(w, http.StatusNoContent, nil)
}
