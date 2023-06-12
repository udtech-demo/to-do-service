package interactor

import (
	"todo-service/src/models"
	"todo-service/src/usecase/presenter"
	"todo-service/src/usecase/repository"
)

type userInteractor struct {
	UserRepository repository.UserRepository
	UserPresenter  presenter.UserPresenter
}

type UserInteractor interface {
	Create(user models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

func NewUserInteractor(
	r repository.UserRepository, p presenter.UserPresenter) UserInteractor {
	return &userInteractor{r, p}
}

func (ui *userInteractor) Create(user models.User) (*models.User, error) {
	return ui.UserRepository.Create(user)
}

func (ui *userInteractor) GetByEmail(email string) (*models.User, error) {
	return ui.UserRepository.GetByEmail(email)
}

func (ui *userInteractor) GetByID(id string) (*models.User, error) {
	return ui.UserRepository.GetByID(id)
}
