package repository

import (
	"product-controller/config"
	"product-controller/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		DB: config.DB,
	}
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	result := r.DB.Create(product)
	return result.Error
}

func (r *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := r.DB.First(&product, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := r.DB.Find(&products)

	return products, result.Error
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	result := r.DB.Save(product)
	return result.Error
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	result := r.DB.Delete(&models.Product{}, id)
	return result.Error
}

func (r *ProductRepository) SaveProductHistory(history *models.ProductHistory) error {
	result := r.DB.Create(history)
	return result.Error
}

func (r *ProductRepository) GetProductHistory(productID uint) ([]models.ProductHistory, error) {
	var history []models.ProductHistory
	result := r.DB.Where("product_id = ?", productID).Find(&history)
	return history, result.Error
}

func (r *ProductRepository) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	result := r.DB.Where("name = ?", name).First(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}
