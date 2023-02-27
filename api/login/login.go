package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go-chi-api/internal/auth"
	"go-chi-api/internal/database"
	"go-chi-api/internal/response"
	"go-chi-api/models"

	"github.com/go-chi/chi/v5"
)

type LoginResource struct{}

func (rs LoginResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Login)
	return r
}

func (rs LoginResource) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	err := user.Validate("login")
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	retrievedUser := models.User{}
	err = database.DB.Debug().Model(models.User{}).Where("username = ?", user.Username).Take(&retrievedUser).Error
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	pwMatch := models.VerifyPassword(retrievedUser.Password, user.Password)

	if !pwMatch {
		response.ERROR(w, http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	token, _ := auth.CreateToken(&retrievedUser)
	fmt.Printf("token: %s", token)
	response.JSON(w, http.StatusOK, token)
}
