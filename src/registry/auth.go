package registry

import (
	interfacePresenter "todo-service/src/interface/presenter"
	interfaceRepository "todo-service/src/interface/repository"
	usecaseInteractor "todo-service/src/usecase/interactor"
	usecasePresenter "todo-service/src/usecase/presenter"
	usecaseRepository "todo-service/src/usecase/repository"
)

func (r *registry) NewAuthInteractor() usecaseInteractor.AuthInteractor {
	return usecaseInteractor.NewAuthInteractor(r.NewAuthRepository(), r.NewUserRepository(), r.jwtConf)
}

func (r *registry) NewAuthMiddleware() usecaseInteractor.Middleware {
	return usecaseInteractor.NewAuthMiddleware(r.NewAuthInteractor(), r.jwtConf)
}

func (r *registry) NewAuthRepository() usecaseRepository.AuthRepository {
	return interfaceRepository.NewAuthRepository(r.db)
}

func (r *registry) NewAuthPresenter() usecasePresenter.AuthPresenter {
	return interfacePresenter.NewAuthPresenter()
}
