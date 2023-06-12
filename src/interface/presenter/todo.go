package presenter

type todoPresenter struct {
}

type TodoPresenter interface {
}

func NewTodoPresenter() TodoPresenter {
	return &todoPresenter{}
}
