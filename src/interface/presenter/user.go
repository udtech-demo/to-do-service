package presenter

type userPresenter struct {
}

type UserPresenter interface {
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}
