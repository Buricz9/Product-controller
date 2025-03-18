package models

import "time"

type ProductHistory struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null"`
	Field     string    `gorm:"size:50;not null"` // Price, Name, Quantity, Description
	OldValue  string    `gorm:"not null"`
	NewValue  string    `gorm:"not null"`
	ChangedAt time.Time `gorm:"autoCreateTime"`
}
