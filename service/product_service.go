package service

import (
	"errors"
	"fmt"
	"product-controller/models"
	"product-controller/repository"
	"regexp"
	"strings"
)

type ProductService struct {
	ProductRepo   *repository.ProductRepository
	BlacklistRepo *repository.BlacklistRepository
}

func NewProductService(productRepo *repository.ProductRepository, blacklistRepo *repository.BlacklistRepository) *ProductService {
	return &ProductService{
		ProductRepo:   productRepo,
		BlacklistRepo: blacklistRepo,
	}
}
func (s *ProductService) AddProduct(product *models.Product) error {
	// Pobierz czarną listę słów
	blacklist, err := s.BlacklistRepo.GetAllBlacklistWords()
	if err != nil {
		return err
	}
	if err = s.validateProduct(product); err != nil {
		return err
	}

	// Sprawdź, czy nazwa produktu zawiera zabronione słowo
	for _, word := range blacklist {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(word.Word)) {
			return errors.New("nazwa produktu zawiera zabronione słowo: " + word.Word)
		}
	}

	// Dodaj produkt
	return s.ProductRepo.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(id uint, updatedProduct *models.Product) error {
	// Pobierz istniejący produkt
	existingProduct, err := s.ProductRepo.GetProductByID(id)
	if err != nil {
		return err
	}

	updatedProduct.ID = id

	if err = s.validateProduct(updatedProduct); err != nil {
		return err
	}

	// Walidacja nazwy z blacklistą
	blacklist, err := s.BlacklistRepo.GetAllBlacklistWords()
	if err != nil {
		return err
	}

	for _, word := range blacklist {
		if strings.Contains(strings.ToLower(updatedProduct.Name), strings.ToLower(word.Word)) {
			return errors.New("nazwa produktu zawiera zabronione słowo: " + word.Word)
		}
	}

	// Zapis historii zmian
	if existingProduct.Name != updatedProduct.Name {
		s.saveProductHistory(id, "Name", existingProduct.Name, updatedProduct.Name)
	}
	if existingProduct.Price != updatedProduct.Price {
		s.saveProductHistory(id, "Price", fmt.Sprintf("%.2f", existingProduct.Price), fmt.Sprintf("%.2f", updatedProduct.Price))
	}
	if existingProduct.Quantity != updatedProduct.Quantity {
		s.saveProductHistory(id, "Quantity", fmt.Sprintf("%d", existingProduct.Quantity), fmt.Sprintf("%d", updatedProduct.Quantity))
	}
	if existingProduct.Description != updatedProduct.Description {
		s.saveProductHistory(id, "Description", existingProduct.Description, updatedProduct.Description)
	}

	// Aktualizacja produktu
	existingProduct.Name = updatedProduct.Name
	existingProduct.Description = updatedProduct.Description
	existingProduct.Price = updatedProduct.Price
	existingProduct.Quantity = updatedProduct.Quantity

	return s.ProductRepo.UpdateProduct(existingProduct)
}

func (s *ProductService) saveProductHistory(productID uint, field, oldValue, newValue string) {
	history := models.ProductHistory{
		ProductID: productID,
		Field:     field,
		OldValue:  oldValue,
		NewValue:  newValue,
	}
	s.ProductRepo.SaveProductHistory(&history)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.ProductRepo.DeleteProduct(id)
}
func (s *ProductService) GetProductHistory(productID uint) ([]models.ProductHistory, error) {
	return s.ProductRepo.GetProductHistory(productID)
}

func (s *ProductService) validateProduct(product *models.Product) error {
	// Walidacja nazwy
	if len(product.Name) < 3 || len(product.Name) > 20 {
		return errors.New("nazwa produktu musi mieć od 3 do 20 znaków")
	}

	namePattern := `^[a-zA-Z0-9]+$`
	matched, _ := regexp.MatchString(namePattern, product.Name)
	if !matched {
		return errors.New("nazwa produktu może zawierać tylko litery i cyfry")
	}

	existing, _ := s.ProductRepo.GetProductByName(product.Name)
	if existing != nil && existing.ID != product.ID {
		return errors.New("produkt o tej nazwie już istnieje")
	}

	var minPrice, maxPrice float64
	switch strings.ToLower(product.Category) {
	case "elektronika":
		minPrice = 50
		maxPrice = 50000
	case "książki":
		minPrice = 5
		maxPrice = 500
	case "odzież":
		minPrice = 10
		maxPrice = 5000
	default:
		return errors.New("kategoria musi być jedną z: Elektronika, Książki, Odzież")
	}

	if product.Price < minPrice || product.Price > maxPrice {
		return fmt.Errorf("cena produktu w kategorii %s musi być w przedziale %.2f - %.2f", product.Category, minPrice, maxPrice)
	}

	if product.Quantity < 0 {
		return errors.New("ilość produktów nie może być ujemna")
	}

	return nil
}
