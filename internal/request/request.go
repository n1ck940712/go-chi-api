package request

import (
	"fmt"
	"go-chi-api/internal/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetParamID(w http.ResponseWriter, r *http.Request) (int, error) {
	customerID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(customerID)
	if err != nil {
		response.ERROR(w, http.StatusNotFound, fmt.Errorf("id param is an invalid integer"))
		return 0, err
	}
	return id, nil
}
