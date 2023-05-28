package customer

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-chi-api/internal/database"
	"go-chi-api/internal/middlewares"
	"go-chi-api/internal/models"
	"go-chi-api/internal/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type CustomerResource struct{}

func (c CustomerResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.BaseAuthentication)

	r.Get("/", GetAllCustomers)
	r.Post("/", CreateCustomer)
	r.Get("/{id}/", GetCustomer)
	r.Put("/{id}/", UpdateCustomer)
	r.Delete("/{id}/", DeleteCustomer)
	return r
}

// Get all customers
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []models.CustomerTable
	database.DB.Find(&customers)
	response.JSON(w, http.StatusOK, customers)
}

// Create a new customer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.CustomerTable
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(customer)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}

		var validationErrorMessages []string
		for _, validationError := range validationErrors {
			validationErrorMessage := fmt.Sprintf("%s validation failed on %s", validationError.Tag(), validationError.Field())
			validationErrorMessages = append(validationErrorMessages, validationErrorMessage)
		}

		response.ERROR(w, http.StatusBadRequest, errors.New(strings.Join(validationErrorMessages, ", ")))
		return
	}

	result := database.DB.Create(&customer)
	if result.Error != nil {
		response.ERROR(w, http.StatusBadRequest, result.Error)
		return
	}
	response.JSON(w, http.StatusCreated, customer)
}

// get single customer by ID
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerID, err := getCustomerID(w, r)
	if err != nil {
		return
	}
	var customer models.CustomerTable
	result := database.DB.First(&customer, customerID)
	if result.Error != nil {
		response.ERROR(w, http.StatusNotFound, result.Error)
		return
	}
	response.JSON(w, http.StatusOK, customer)
}

// update customer
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customerID, err1 := getCustomerID(w, r)
	if err1 != nil {
		return
	}
	var customer models.CustomerTable
	result := database.DB.First(&customer, customerID)
	if result.Error != nil {
		response.ERROR(w, http.StatusNotFound, result.Error)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(customer)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}

		var validationErrorMessages []string
		for _, validationError := range validationErrors {
			validationErrorMessage := fmt.Sprintf("%s validation failed on %s", validationError.Tag(), validationError.Field())
			validationErrorMessages = append(validationErrorMessages, validationErrorMessage)
		}

		response.ERROR(w, http.StatusBadRequest, errors.New(strings.Join(validationErrorMessages, ", ")))
		return
	}

	result = database.DB.Save(&customer)
	if result.Error != nil {
		response.ERROR(w, http.StatusBadRequest, result.Error)
		return
	}
	response.JSON(w, http.StatusOK, customer)
}

// delete customer
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerID, err := getCustomerID(w, r)
	if err != nil {
		return
	}

	rowsDeleted := database.DB.Delete(&models.CustomerTable{}, customerID).RowsAffected
	if rowsDeleted == 0 {
		response.ERROR(w, http.StatusNotFound, fmt.Errorf("customer with id %d not found", customerID))
		return
	}

}

// check if customerID param is a valid integer
func getCustomerID(w http.ResponseWriter, r *http.Request) (int, error) {
	customerID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(customerID)
	if err != nil {
		response.ERROR(w, http.StatusNotFound, fmt.Errorf("id param is an invalid integer"))
		return 0, err
	}
	return id, nil
}
