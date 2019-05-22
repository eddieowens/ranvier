package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/server/exchange/response"
	"net/http"
)

const ConfigControllerServiceKey = "ConfigControllerService"

type ConfigControllerService interface {
	Query(name string, query string) (resp response.Config, err error)
}

type configControllerServiceImpl struct {
	MappingService     MappingService     `inject:"MappingService"`
	ConfigQueryService ConfigQueryService `inject:"ConfigQueryService"`
}

func (g *configControllerServiceImpl) Query(name string, query string) (resp response.Config, err error) {
	config, err := g.ConfigQueryService.Query(name, query)
	if err != nil {
		return resp, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid query: %s", err.Error()))
	}
	if config != nil {
		return resp, NewKeyNotFoundError(query)
	}

	return g.MappingService.ToConfig(config), nil
}
