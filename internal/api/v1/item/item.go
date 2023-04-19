package item

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"

	"github.com/go-chi/chi/v5"
)

type ItemsResource struct{}

func (rs ItemsResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.BaseAuthentication)
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
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
	}
	json.Unmarshal(reqBody, &item)

	var request struct {
		ItemTypeID int64 `json:"item_type_id"`
	}
	json.Unmarshal(reqBody, &request)
	var itemType models.ItemTypeTable
	if err := database.DB.First(&itemType, request.ItemTypeID).Error; err != nil {
		fmt.Println(err)
	}
	item.ItemType = itemType
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
	idStr := r.Context().Value("id").(string)
	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	result := database.DB.Where("id=?", idInt).Find(&models.ItemTable{})
	if result.RowsAffected == 0 {
		response.ERROR(w, http.StatusNotFound, fmt.Errorf("item not found"))
		return
	}
	item := models.ItemTable{}
	json.NewDecoder(r.Body).Decode(&item)
	database.DB.Where("id=?", idInt).Updates(item)

	updated_item := models.ItemTable{}
	database.DB.Find(&updated_item, "id=?", idInt)
	response.JSON(w, http.StatusOK, updated_item)
}

func (rs ItemsResource) Delete(w http.ResponseWriter, r *http.Request) {
	item := models.ItemTable{}
	id := r.Context().Value("id").(string)
	result := database.DB.Delete(&item, id)
	if result.RowsAffected == 0 {
		response.ERROR(w, http.StatusNotFound, fmt.Errorf("item not found"))
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func PrintAttributes(i interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("%s: %v\n", field.Name, value)
	}
}
