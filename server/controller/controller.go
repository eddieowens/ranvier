package controller

import (
	"github.com/two-rabbits/ranvier/server/model"
)

type Controller interface {
	GetRoutes() []model.Route
}
