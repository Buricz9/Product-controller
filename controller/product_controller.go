package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"product-controller/models"
	"product-controller/service"
	"strconv"
	"strings"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (c *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.ProductService.ProductRepo.GetAllProducts()
	if err != nil {
		http.Error(w, "Błąd pobierania produktów", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Niepoprawne dane wejściowe", http.StatusBadRequest)
		return
	}

	err = c.ProductService.AddProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (c *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID produktu", http.StatusBadRequest)
		return
	}

	var updatedProduct models.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Niepoprawne dane wejściowe", http.StatusBadRequest)
		return
	}

	err = c.ProductService.UpdateProduct(uint(id), &updatedProduct)
	if err != nil {
		if strings.Contains(err.Error(), "nie istnieje") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}

func (c *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID produktu", http.StatusBadRequest)
		return
	}

	err = c.ProductService.DeleteProduct(uint(id))
	if err != nil {
		http.Error(w, "Błąd usuwania produktu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (c *ProductController) GetProductHistory(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID produktu", http.StatusBadRequest)
		return
	}

	history, err := c.ProductService.GetProductHistory(uint(id))
	if err != nil {
		http.Error(w, "Błąd pobierania historii produktu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (c *ProductController) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID produktu", http.StatusBadRequest)
		return
	}

	product, err := c.ProductService.ProductRepo.GetProductByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Produkt nie znaleziony", http.StatusNotFound)
		} else {
			http.Error(w, "Błąd pobierania produktu: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
