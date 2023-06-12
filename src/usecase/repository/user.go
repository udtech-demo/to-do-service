package repository

import (
	"todo-service/src/models"
)

type UserRepository interface {
	Create(user models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}
