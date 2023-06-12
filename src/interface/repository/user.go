package repository

import (
	"strings"
	"todo-service/src/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(user models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(user models.User) (*models.User, error) {

	if err := ur.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetByEmail(email string) (*models.User, error) {

	var user models.User
	if err := ur.db.Model(user).Where("email = ?", strings.ToLower(email)).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetByID(id string) (*models.User, error) {

	var user models.User
	if err := ur.db.Model(user).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
