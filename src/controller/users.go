package controller

type UsersController interface {
	Controller
}

type usersControllerImpl struct {
}

func (u *usersControllerImpl) GetRoutes() []interface{} {
}
