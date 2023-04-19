package restock

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RestockResource struct{}

func (rs RestockResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.BaseAuthentication)
	r.Get("/", GetAllRestocks)
	r.Post("/", CreateRestock)
	r.Get("/{restockID}/", GetRestock)
	// r.Put("/{restockID}", UpdateRestock)
	// r.Delete("/{restockID}", DeleteRestock)

	return r
}

func GetAllRestocks(w http.ResponseWriter, r *http.Request) {
	var restocks []models.RestockTable
	database.DB.Preload("RestockItems").Find(&restocks)
	response.JSON(w, http.StatusOK, restocks)
}

func GetRestock(w http.ResponseWriter, r *http.Request) {
	restockID := chi.URLParam(r, "restockID")
	var restock models.RestockTable
	result := database.DB.Preload("RestockItems").First(&restock, restockID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.ERROR(w, http.StatusNotFound, result.Error)
			return
		}
		response.ERROR(w, http.StatusInternalServerError, result.Error)
		return
	}
	response.JSON(w, http.StatusOK, restock)
}

func CreateRestock(w http.ResponseWriter, r *http.Request) {
	request := RestockItemsRequest{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
	}
	json.Unmarshal(reqBody, &request)
	if err := request.Validate(); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("user").(*models.User)

	var restockTable models.RestockTable
	json.Unmarshal(reqBody, &restockTable)

	restockTable.UserID = user.ID

	database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Create(&restockTable).Error; err != nil {
			return err
		}
		for _, item := range request.Items {
			restockItem := models.RestockItemTable{
				Restock:    restockTable,
				ItemID:     item.ItemID,
				Quantity:   item.Quantity,
				UnitPrice:  item.Price,
				TotalPrice: decimal.NewFromFloat(float64(item.Quantity)).Mul(item.Price),
			}
			if err := database.DB.Create(&restockItem).Error; err != nil {
				return err
			}
			database.DB.Model(&models.ItemTable{}).Where("id = ?", item.ItemID).Update("quantity", gorm.Expr("quantity + ?", item.Quantity))
		}
		return nil
	})
	database.DB.Preload("RestockItems").Find(&restockTable, restockTable.ID)
	response.JSON(w, http.StatusOK, restockTable)
}

type RestockItemsRequest struct {
	Description string               `json:"description"`
	Items       []RestockItemRequest `json:"items" validate:"required"`
}

type RestockItemRequest struct {
	ItemID   int32           `json:"item_id" validate:"required,gte=1"`
	Quantity int32           `json:"quantity" validate:"required,gte=1"`
	Price    decimal.Decimal `json:"price" validate:"required,gte=0"`
}

func (r *RestockItemsRequest) Validate() error {
	// Validate request
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}

	// Check if item exists
	for _, item := range r.Items {
		var itemTable models.ItemTable
		result := database.DB.First(&itemTable, item.ItemID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("item with id %d does not exist", item.ItemID)
			}
			return result.Error
		}
	}

	return nil
}
