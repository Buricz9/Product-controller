package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"product-controller/config"
	"product-controller/controller"
	"product-controller/models"
	"product-controller/repository"
	"product-controller/service"
	"strconv"
	"testing"
)

func setupRouter() http.Handler {
	config.InitDB()

	productRepo := repository.NewProductRepository()
	blacklistRepo := repository.NewBlacklistRepository()

	productService := service.NewProductService(productRepo, blacklistRepo)

	productController := controller.NewProductController(productService)
	blacklistController := controller.NewBlacklistController(blacklistRepo)

	r := chi.NewRouter()

	// Product routes
	r.Get("/products", productController.GetAllProducts)
	r.Get("/products/{id}", productController.GetProductByID)
	r.Post("/products", productController.AddProduct)
	r.Put("/products/{id}", productController.UpdateProduct)
	r.Delete("/products/{id}", productController.DeleteProduct)
	r.Get("/products/{id}/history", productController.GetProductHistory)

	// Blacklist routes
	r.Get("/blacklist", blacklistController.GetAllBlacklistWords)
	r.Post("/blacklist", blacklistController.AddBlacklistWord)
	r.Delete("/blacklist/{id}", blacklistController.DeleteBlacklistWord)

	return r
}
func TestGetAllProducts(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
func TestCreateProduct(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "Test Product",
		Description: "Description",
		Price:       100.0,
		Quantity:    10,
	}

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdProduct models.Product
	json.Unmarshal(rr.Body.Bytes(), &createdProduct)
	assert.Equal(t, product.Name, createdProduct.Name)
}
func TestGetProductByID(t *testing.T) {
	router := setupRouter()

	// Najpierw dodaj produkt
	product := models.Product{
		Name:        "Test Product 2",
		Description: "Another Description",
		Price:       200.0,
		Quantity:    5,
	}
	body, _ := json.Marshal(product)
	reqCreate, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()
	router.ServeHTTP(rrCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var createdProduct models.Product
	json.Unmarshal(rrCreate.Body.Bytes(), &createdProduct)

	// Pobierz po ID
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(createdProduct.ID)), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
func TestDeleteProduct(t *testing.T) {
	router := setupRouter()

	// Najpierw dodaj produkt
	product := models.Product{
		Name:        "Test Product 3",
		Description: "Delete Me",
		Price:       150.0,
		Quantity:    3,
	}
	body, _ := json.Marshal(product)
	reqCreate, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()
	router.ServeHTTP(rrCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var createdProduct models.Product
	json.Unmarshal(rrCreate.Body.Bytes(), &createdProduct)

	// Usuń produkt
	req, _ := http.NewRequest("DELETE", "/products/"+strconv.Itoa(int(createdProduct.ID)), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
