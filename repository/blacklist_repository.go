package repository

import (
	"product-controller/config"
	"product-controller/models"

	"gorm.io/gorm"
)

type BlacklistRepository struct {
	DB *gorm.DB
}

func NewBlacklistRepository() *BlacklistRepository {
	return &BlacklistRepository{
		DB: config.DB,
	}
}

func (r *BlacklistRepository) GetAllBlacklistWords() ([]models.BlacklistWord, error) {
	var words []models.BlacklistWord
	result := r.DB.Find(&words)
	return words, result.Error
}

func (r *BlacklistRepository) AddBlacklistWord(word *models.BlacklistWord) error {
	result := r.DB.Create(word)
	return result.Error
}

func (r *BlacklistRepository) DeleteBlacklistWord(id uint) error {
	result := r.DB.Delete(&models.BlacklistWord{}, id)
	return result.Error
}
