package service

import (
	"github.com/eddieowens/ranvier/server/app/exchange/response"
)

const ConfigControllerServiceKey = "ConfigControllerService"

type ConfigControllerService interface {
	Query(name string, query string) (resp response.Config, err error)
}

type configControllerServiceImpl struct {
	MappingService MappingService `inject:"MappingService"`
	ConfigService  ConfigService  `inject:"ConfigService"`
}

func (g *configControllerServiceImpl) Query(name string, query string) (resp response.Config, err error) {
	config, err := g.ConfigService.Query(name, query)
	if err != nil {
		return resp, err
	}

	return g.MappingService.ToConfig(config), nil
}
