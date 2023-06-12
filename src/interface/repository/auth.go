package repository

import (
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}
