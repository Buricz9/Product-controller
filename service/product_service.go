package service

import (
	"errors"
	"product-controller/models"
	"product-controller/repository"
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

	// Sprawdź, czy nazwa produktu zawiera zabronione słowo
	for _, word := range blacklist {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(word.Word)) {
			return errors.New("nazwa produktu zawiera zabronione słowo: " + word.Word)
		}
	}

	// Dodaj produkt
	return s.ProductRepo.CreateProduct(product)
}
