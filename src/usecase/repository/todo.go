package repository

import (
	"todo-service/graph/model"
	"todo-service/src/models"
)

type TodoRepository interface {
	Create(input model.NewTodo, userId string) (*models.Todo, error)
	MarkComplete(id string, userId string) error
	Delete(id string, userId string) (bool, error)
	GetByID(id string, userId string) (*models.Todo, error)
	List(userId string) ([]*models.Todo, error)
}
