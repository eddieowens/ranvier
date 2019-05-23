package controller

import (
	"github.com/two-rabbits/ranvier/server/app/model"
)

type Controller interface {
	GetRoutes() []model.Route
}
