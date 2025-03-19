package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-controller/config"
	"product-controller/controller"
	"product-controller/models"
	"product-controller/repository"
	"product-controller/service"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupRouter() http.Handler {
	config.InitDB()
	truncateTables()

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

/////////////////////////////////////////////////////
//                   Produkty                     //
/////////////////////////////////////////////////////

func TestGetAllProducts(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateProductSuccess(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "TestProduct",
		Category:    "Elektronika",
		Description: "Opis produktu",
		Price:       1000.0,
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

	product := models.Product{
		Name:        "GetByIDProduct",
		Category:    "Elektronika",
		Description: "Opis produktu",
		Price:       500.0,
		Quantity:    2,
	}
	body, _ := json.Marshal(product)
	reqCreate, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()

	router.ServeHTTP(rrCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var createdProduct models.Product
	json.Unmarshal(rrCreate.Body.Bytes(), &createdProduct)

	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(createdProduct.ID)), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateProductAndCheckHistory(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "UpdateTestProduct",
		Category:    "Elektronika",
		Description: "Opis produktu",
		Price:       500.0,
		Quantity:    2,
	}
	body, _ := json.Marshal(product)
	reqCreate, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()

	router.ServeHTTP(rrCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var createdProduct models.Product
	json.Unmarshal(rrCreate.Body.Bytes(), &createdProduct)

	updatedProduct := models.Product{
		Name:        "UpdatedProductName",
		Category:    "Elektronika",
		Description: "Nowy opis",
		Price:       750.0,
		Quantity:    5,
	}
	bodyUpdate, _ := json.Marshal(updatedProduct)
	reqUpdate, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(createdProduct.ID)), bytes.NewBuffer(bodyUpdate))
	reqUpdate.Header.Set("Content-Type", "application/json")
	rrUpdate := httptest.NewRecorder()

	router.ServeHTTP(rrUpdate, reqUpdate)
	assert.Equal(t, http.StatusOK, rrUpdate.Code)

	reqHistory, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(createdProduct.ID))+"/history", nil)
	rrHistory := httptest.NewRecorder()

	router.ServeHTTP(rrHistory, reqHistory)
	assert.Equal(t, http.StatusOK, rrHistory.Code)

	var history []models.ProductHistory
	json.Unmarshal(rrHistory.Body.Bytes(), &history)

	assert.True(t, len(history) > 0)
}

func TestDeleteProduct(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "DeleteMeProduct",
		Category:    "Elektronika",
		Description: "Opis produktu",
		Price:       300.0,
		Quantity:    1,
	}
	body, _ := json.Marshal(product)
	reqCreate, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()

	router.ServeHTTP(rrCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, rrCreate.Code)

	var createdProduct models.Product
	json.Unmarshal(rrCreate.Body.Bytes(), &createdProduct)

	reqDelete, _ := http.NewRequest("DELETE", "/products/"+strconv.Itoa(int(createdProduct.ID)), nil)
	rrDelete := httptest.NewRecorder()

	router.ServeHTTP(rrDelete, reqDelete)
	assert.Equal(t, http.StatusNoContent, rrDelete.Code)
}

/////////////////////////////////////////////////////
//                 Walidacje                      //
/////////////////////////////////////////////////////

func TestCreateProductWithShortName(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "ab",
		Category:    "Elektronika",
		Description: "Za krótka nazwa",
		Price:       1000.0,
		Quantity:    1,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "nazwa produktu musi mieć od 3 do 20 znaków")
}

func TestCreateProductWithInvalidCategory(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "ValidName",
		Category:    "Inne",
		Description: "Zła kategoria",
		Price:       100.0,
		Quantity:    1,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "kategoria musi być jedną z")
}

func TestCreateProductWithNegativeQuantity(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "NegativeQuantP",
		Category:    "Elektronika",
		Description: "Ujemna ilość",
		Price:       1000.0,
		Quantity:    -5,
	}

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "ilość produktów nie może być ujemna")
}

func TestCreateProductWithInvalidPrice(t *testing.T) {
	router := setupRouter()

	product := models.Product{
		Name:        "InvalidPriceProduct",
		Category:    "Elektronika",
		Description: "Za mała cena",
		Price:       10.0, // Minimum dla elektroniki to 50
		Quantity:    1,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "cena produktu w kategorii")
}

/////////////////////////////////////////////////////
//                  Blacklist                     //
/////////////////////////////////////////////////////

func TestAddBlacklistWord(t *testing.T) {
	router := setupRouter()

	word := models.BlacklistWord{
		Word: "zabronione",
	}
	body, _ := json.Marshal(word)
	req, _ := http.NewRequest("POST", "/blacklist", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdWord models.BlacklistWord
	json.Unmarshal(rr.Body.Bytes(), &createdWord)
	assert.Equal(t, word.Word, createdWord.Word)
}

func TestGetBlacklist(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/blacklist", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var blacklist []models.BlacklistWord
	err := json.Unmarshal(rr.Body.Bytes(), &blacklist)
	assert.NoError(t, err)
}

func TestDeleteBlacklistWord(t *testing.T) {
	router := setupRouter()

	word := models.BlacklistWord{
		Word: "do_usuniecia",
	}
	body, _ := json.Marshal(word)
	reqAdd, _ := http.NewRequest("POST", "/blacklist", bytes.NewBuffer(body))
	reqAdd.Header.Set("Content-Type", "application/json")
	rrAdd := httptest.NewRecorder()

	router.ServeHTTP(rrAdd, reqAdd)
	assert.Equal(t, http.StatusCreated, rrAdd.Code)

	var createdWord models.BlacklistWord
	json.Unmarshal(rrAdd.Body.Bytes(), &createdWord)

	reqDelete, _ := http.NewRequest("DELETE", "/blacklist/"+strconv.Itoa(int(createdWord.ID)), nil)
	rrDelete := httptest.NewRecorder()

	router.ServeHTTP(rrDelete, reqDelete)
	assert.Equal(t, http.StatusNoContent, rrDelete.Code)
}

func TestCreateProductWithBlacklistedWord(t *testing.T) {
	router := setupRouter()

	word := models.BlacklistWord{
		Word: "blokuj",
	}
	bodyWord, _ := json.Marshal(word)
	reqWord, _ := http.NewRequest("POST", "/blacklist", bytes.NewBuffer(bodyWord))
	reqWord.Header.Set("Content-Type", "application/json")
	rrWord := httptest.NewRecorder()

	router.ServeHTTP(rrWord, reqWord)
	assert.Equal(t, http.StatusCreated, rrWord.Code)

	product := models.Product{
		Name:        "SuperBlokujLaptop",
		Category:    "Elektronika",
		Description: "Test",
		Price:       1000.0,
		Quantity:    10,
	}
	bodyProduct, _ := json.Marshal(product)
	reqProduct, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(bodyProduct))
	reqProduct.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, reqProduct)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "zawiera zabronione słowo")
}

func truncateTables() {
	db := config.DB

	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	db.Exec("TRUNCATE TABLE products;")
	db.Exec("TRUNCATE TABLE product_histories;")
	db.Exec("TRUNCATE TABLE blacklist_words;")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
}
