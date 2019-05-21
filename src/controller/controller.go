package controller

import (
	"github.com/two-rabbits/ranvier/src/model"
)

type Controller interface {
	GetRoutes() []model.Route
}
