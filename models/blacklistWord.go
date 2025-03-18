package models

type BlacklistWord struct {
	ID   uint   `gorm:"primaryKey"`
	Word string `gorm:"size:255;not null;unique"`
}
