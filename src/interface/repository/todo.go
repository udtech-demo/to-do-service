package repository

import (
	"time"
	"todo-service/graph/model"
	"todo-service/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

type TodoRepository interface {
	Create(input model.NewTodo, userId string) (*models.Todo, error)
	MarkComplete(id string, userId string) error
	Delete(id string, userId string) (bool, error)
	GetByID(id string, userId string) (*models.Todo, error)
	List(userId string) ([]*models.Todo, error)
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db}
}

func (ur *todoRepository) Create(input model.NewTodo, userId string) (*models.Todo, error) {

	todo := models.Todo{
		ID:      uuid.New(),
		Text:    input.Text,
		UserID:  uuid.MustParse(userId),
		Done:    false,
		Created: time.Now(),
	}

	if err := ur.db.Model(todo).Create(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (ur *todoRepository) MarkComplete(id string, userId string) error {

	if err := ur.db.Model((*models.Todo)(nil)).Where("id = ? AND user_id = ?", id, userId).Update("done", true).Error; err != nil {
		return err
	}

	return nil
}

func (ur *todoRepository) Delete(id string, userId string) (bool, error) {
	var todo models.Todo
	q := ur.db.Where("id = ? AND user_id = ?", id, userId).Delete(&todo)
	if q.Error != nil {
		return false, q.Error
	}

	if q.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (ur *todoRepository) GetByID(id string, userId string) (*models.Todo, error) {

	var todo models.Todo
	if err := ur.db.Model(todo).Where("id = ? AND user_id = ?", id, userId).Preload("User").Take(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (ur *todoRepository) List(userId string) ([]*models.Todo, error) {

	var todos []*models.Todo
	if err := ur.db.Model(todos).Where("user_id = ?", userId).Preload("User").Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}
