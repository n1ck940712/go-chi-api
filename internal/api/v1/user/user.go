package user

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"

	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"

	"github.com/go-chi/chi/v5"
)

type UsersResource struct{}

func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.Create)

	r.With(middlewares.IsAdmin).Get("/", rs.List)

	r.With(middlewares.IsAdmin).Route("/{id}", func(r chi.Router) {
		r.Use(ctx)
		r.Get("/", rs.Get)
		// r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	users := models.User{}
	if err := database.DB.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, users)

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
	user := models.User{}
	id := r.Context().Value("id").(string)
	idStr, _ := strconv.Atoi(id)

	if retrievedUser, err := user.Get(database.DB, int32(idStr)); err != nil {
		response.JSON(w, http.StatusNotFound, err.Error())
	} else {
		response.JSON(w, http.StatusOK, retrievedUser)
	}

}

// func (rs UsersResource) Update(w http.ResponseWriter, r *http.Request) {
// 	id := r.Context().Value("id").(string)
// 	User := models.User{}
// 	json.NewDecoder(r.Body).Decode(&User)
// 	database.DB.Model(&User).Where("id = ?", id).UpdateColumns(
// 		map[string]interface{}{
// 			"username":   User.Username,
// 			"password":   User.Password,
// 			"updated_at": time.Now(),
// 		},
// 	)
// 	updated_User := models.User{}
// 	database.DB.Find(&updated_User, id)
// 	response.JSON(w, http.StatusOK, updated_User)
// }

func (rs UsersResource) Delete(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		response.JSON(w, http.StatusNotFound, "User not found")
		return
	}

	result := database.DB.Delete(&user)

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	response.JSON(w, http.StatusOK, fmt.Sprintf("User ID (%d) deleted", user.ID))
}

func getUserFromCtx(r *http.Request) *models.User {
	// Retrieve the user ID from the request context
	userID := r.Context().Value("id")
	// Query the database for the user with the given ID
	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		fmt.Println("Could not find user with ID", userID)
		return nil
	}

	return &user
}
