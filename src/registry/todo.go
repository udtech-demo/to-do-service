package registry

import (
	interfacePresenter "todo-service/src/interface/presenter"
	interfaceRepository "todo-service/src/interface/repository"
	usecaseInteractor "todo-service/src/usecase/interactor"
	usecasePresenter "todo-service/src/usecase/presenter"
	usecaseRepository "todo-service/src/usecase/repository"
)

func (r *registry) NewTodoInteractor() usecaseInteractor.TodoInteractor {
	return usecaseInteractor.NewTodoInteractor(r.NewTodoRepository(), r.NewTodoPresenter())
}

func (r *registry) NewTodoRepository() usecaseRepository.TodoRepository {
	return interfaceRepository.NewTodoRepository(r.db)
}

func (r *registry) NewTodoPresenter() usecasePresenter.TodoPresenter {
	return interfacePresenter.NewTodoPresenter()
}
