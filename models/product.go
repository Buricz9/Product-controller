package models

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:255;not null;unique"`
	Category    string  `gorm:"size:50;not null"`
	Description string  `gorm:"size:1000"`
	Price       float64 `gorm:"not null"`
	Quantity    int     `gorm:"not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
