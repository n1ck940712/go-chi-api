package item_type

import (
	"encoding/json"
	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type ItemTableInput struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type ItemTypesResource struct{}

func (rs ItemTypesResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.BaseAuthentication)
	r.Get("/", GetAllItemTypes)
	r.Post("/", CreateItemType)
	r.Get("/{itemTypeID}/", GetItemType)
	r.Put("/{itemTypeID}/", UpdateItemType)
	r.Delete("/{itemTypeID}/", DeleteItemType)
	return r
}

// Get all item types
func GetAllItemTypes(w http.ResponseWriter, r *http.Request) {
	var itemTypes []models.ItemTypeTable
	database.DB.Find(&itemTypes)
	response.JSON(w, http.StatusOK, itemTypes)
}

// Get single item type by ID
func GetItemType(w http.ResponseWriter, r *http.Request) {
	itemTypeID := chi.URLParam(r, "itemTypeID")
	var itemType models.ItemTypeTable
	result := database.DB.First(&itemType, itemTypeID)
	if result.Error != nil {
		response.ERROR(w, http.StatusNotFound, result.Error)
		return
	}
	response.JSON(w, http.StatusOK, itemType)
}

// Create a new item type
func CreateItemType(w http.ResponseWriter, r *http.Request) {
	var itemType models.ItemTypeTable
	var input ItemTableInput
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	json.Unmarshal(reqBody, &input)
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	json.Unmarshal(reqBody, &itemType)
	database.DB.Create(&itemType)
	response.JSON(w, http.StatusOK, itemType)
}

// Update an existing item type
func UpdateItemType(w http.ResponseWriter, r *http.Request) {
	var input ItemTableInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	itemTypeID := chi.URLParam(r, "itemTypeID")
	var itemType models.ItemTypeTable
	result := database.DB.First(&itemType, itemTypeID)
	if result.RowsAffected == 0 {
		response.ERROR(w, http.StatusNotFound, result.Error)
		return
	}

	itemType.Name = input.Name
	itemType.Description = input.Description

	database.DB.Save(&itemType)
	response.JSON(w, http.StatusOK, itemType)
}

// Delete an item type
func DeleteItemType(w http.ResponseWriter, r *http.Request) {
	itemTypeID := chi.URLParam(r, "itemTypeID")
	var itemType models.ItemTypeTable
	database.DB.First(&itemType, itemTypeID)
	if itemType.ID != 0 {
		database.DB.Delete(&itemType)
		response.JSON(w, http.StatusOK, "Item type deleted")
	} else {
		response.JSON(w, http.StatusNoContent, nil)
	}
}
