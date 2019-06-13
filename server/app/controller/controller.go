package controller

import (
	"github.com/eddieowens/ranvier/server/app/model"
)

type Controller interface {
	GetRoutes() []model.Route
}
