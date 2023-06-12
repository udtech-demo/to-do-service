package registry

import (
	interfacePresenter "todo-service/src/interface/presenter"
	interfaceRepository "todo-service/src/interface/repository"
	usecaseInteractor "todo-service/src/usecase/interactor"
	usecasePresenter "todo-service/src/usecase/presenter"
	usecaseRepository "todo-service/src/usecase/repository"
)

func (r *registry) NewUserInteractor() usecaseInteractor.UserInteractor {
	return usecaseInteractor.NewUserInteractor(r.NewUserRepository(), r.NewUserPresenter())
}

func (r *registry) NewUserRepository() usecaseRepository.UserRepository {
	return interfaceRepository.NewUserRepository(r.db)
}

func (r *registry) NewUserPresenter() usecasePresenter.UserPresenter {
	return interfacePresenter.NewUserPresenter()
}
