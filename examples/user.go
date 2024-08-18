package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"index"`
	Email string `gorm:"uniqueIndex"`
	Age   int
}