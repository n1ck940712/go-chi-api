package order

import (
	"encoding/json"
	"fmt"
	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/request"
	"go-chi-api/internal/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderResource struct{}

func (c OrderResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.BaseAuthentication)
	r.Get("/", GetAllOrders)
	r.Post("/", CreateOrder)
	r.Get("/{id}/", GetOrder)
	// r.Put("/{id}/", UpdateOrder)
	// r.Delete("/{id}/", DeleteOrder)
	return r
}

// Get all orders
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.OrderTable
	database.DB.Find(&orders)
	response.JSON(w, http.StatusOK, orders)
}

// Get order by id
func GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := request.GetParamID(w, r)
	if err != nil {
		return
	}
	var order models.OrderTable
	result := database.DB.First(&order, orderID)
	if result.Error != nil {
		response.ERROR(w, http.StatusNotFound, result.Error)
		return
	}
	response.JSON(w, http.StatusOK, order)
}

// Create order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.OrderTable
	json.NewDecoder(r.Body).Decode(&order)
	// validate := validator.New()
	// err := validate.Struct(order)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	order.Status = "pending"

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		orderTotal := decimal.NewFromInt(0)
		for _, item := range order.OrderItems {
			retrievedItem, err := models.GetItemFromID(item.ItemID)
			if err != nil {
				return err
			}
			if retrievedItem.Quantity < item.Quantity {
				return fmt.Errorf("not enough %s (item id %d)", retrievedItem.Name, retrievedItem.ID)
			}
			retrievedItem.Quantity -= item.Quantity
			result := tx.Save(&retrievedItem)
			if result.Error != nil {
				response.ERROR(w, http.StatusBadRequest, result.Error)
				return result.Error
			}
			quantityDecimal := decimal.NewFromInt32(item.Quantity)
			totalPrice := retrievedItem.Price.Mul(quantityDecimal)
			orderTotal = orderTotal.Add(totalPrice)
			saveItem := models.OrderItemTable{
				Quantity:   item.Quantity,
				Order:      order,
				UnitPrice:  retrievedItem.Price,
				TotalPrice: totalPrice,
			}
			if updateErr := tx.Model(&item).Where("id = ?", item.ID).Updates(saveItem); updateErr.Error != nil {
				return updateErr.Error
			}

		}
		order.OrderTotal = orderTotal
		if orderUpdate := tx.Model(&order).Where("id = ?", order.ID).Update("order_total", orderTotal); orderUpdate.Error != nil {
			return orderUpdate.Error
		}
		return nil
	})
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response.JSON(w, http.StatusCreated, order)
}
