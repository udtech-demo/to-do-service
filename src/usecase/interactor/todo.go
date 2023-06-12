package interactor

import (
	"todo-service/graph/model"
	"todo-service/src/models"
	"todo-service/src/usecase/presenter"
	"todo-service/src/usecase/repository"
)

type todoInteractor struct {
	TodoRepository repository.TodoRepository
	TodoPresenter  presenter.TodoPresenter
}

type TodoInteractor interface {
	Create(input model.NewTodo, userId string) (*models.Todo, error)
	MarkComplete(id string, userId string) (*models.Todo, error)
	Delete(id string, userId string) (bool, error)
	GetByID(id string, userId string) (*models.Todo, error)
	List(userId string) ([]*models.Todo, error)
}

func NewTodoInteractor(
	r repository.TodoRepository, p presenter.TodoPresenter) TodoInteractor {
	return &todoInteractor{r, p}
}

func (ti *todoInteractor) Create(input model.NewTodo, userId string) (*models.Todo, error) {
	todo, err := ti.TodoRepository.Create(input, userId)
	if err != nil {
		return nil, err
	}

	return ti.TodoRepository.GetByID(todo.ID.String(), userId)
}

func (ti *todoInteractor) MarkComplete(id string, userId string) (*models.Todo, error) {
	err := ti.TodoRepository.MarkComplete(id, userId)
	if err != nil {
		return nil, err
	}

	return ti.TodoRepository.GetByID(id, userId)
}

func (ti *todoInteractor) Delete(id string, userId string) (bool, error) {
	return ti.TodoRepository.Delete(id, userId)
}

func (ti *todoInteractor) GetByID(id string, userId string) (*models.Todo, error) {
	return ti.TodoRepository.GetByID(id, userId)
}

func (ti *todoInteractor) List(userId string) ([]*models.Todo, error) {
	return ti.TodoRepository.List(userId)
}
