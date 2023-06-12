package registry

import (
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/usecase/interactor"

	"gorm.io/gorm"
)

type UseCase struct {
	AuthMiddleware interface{ interactor.Middleware }
	User           interface{ interactor.UserInteractor }
	Todo           interface{ interactor.TodoInteractor }
	Auth           interface{ interactor.AuthInteractor }
}

type registry struct {
	db      *gorm.DB
	jwtConf authentication.JwtConfigurator
}

type Registry interface {
	NewUseCase() UseCase
}

func NewRegistry(db *gorm.DB, jc authentication.JwtConfigurator) Registry {
	return &registry{
		db:      db,
		jwtConf: jc,
	}
}

func (r *registry) NewUseCase() UseCase {
	return UseCase{
		AuthMiddleware: r.NewAuthMiddleware(),
		User:           r.NewUserInteractor(),
		Todo:           r.NewTodoInteractor(),
		Auth:           r.NewAuthInteractor(),
	}
}
