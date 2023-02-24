package user

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"time"

	"go-chi-api/internal/database"
	"go-chi-api/internal/response"
	"go-chi-api/models"

	"github.com/go-chi/chi/v5"
)

type UsersResource struct{}

func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

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

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	result := []models.User{}
	database.DB.Find(&result)
	response.JSON(w, http.StatusOK, result)
}

func (rs UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	User := models.User{}
	json.NewDecoder(r.Body).Decode(&User)
	err := User.Validate("")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
	}
	userCreated, err := User.Save(database.DB)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
	}
	response.JSON(w, http.StatusOK, userCreated)
}

func ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs UsersResource) Get(w http.ResponseWriter, r *http.Request) {
	result := models.User{}
	id := r.Context().Value("id").(string)
	database.DB.Find(&result, id)
	response.JSON(w, http.StatusOK, result)

}

func (rs UsersResource) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	User := models.User{}
	json.NewDecoder(r.Body).Decode(&User)
	database.DB.Model(&User).Where("id = ?", id).UpdateColumns(
		map[string]interface{}{
			"username":   User.Username,
			"password":   User.Password,
			"updated_at": time.Now(),
		},
	)
	updated_User := models.User{}
	database.DB.Find(&updated_User, id)
	response.JSON(w, http.StatusOK, updated_User)
}

func (rs UsersResource) Delete(w http.ResponseWriter, r *http.Request) {
	User := models.User{}
	id := r.Context().Value("id").(string)
	database.DB.Delete(&User, id)
	response.JSON(w, http.StatusNoContent, nil)
}
