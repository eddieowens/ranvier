package service

import (
	"fmt"
	"github.com/eddieowens/ranvier/server/app/exchange/response"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/labstack/echo"
)

const ConfigControllerServiceKey = "ConfigControllerService"

type ConfigControllerService interface {
	Query(name string, query string) (*response.Config, error)
	Delete(name string) (*response.Config, error)
	Create(conf *model.Config) (*response.Config, error)
	Update(config *model.Config) (*response.Config, error)
}

type configControllerServiceImpl struct {
	MappingService MappingService `inject:"MappingService"`
	ConfigService  ConfigService  `inject:"ConfigService"`
}

func (c *configControllerServiceImpl) Delete(name string) (*response.Config, error) {
	conf := c.ConfigService.Delete(name)
	if conf == nil {
		return nil, echo.NewHTTPError(400, fmt.Sprintf("%s could not be found", name))
	}

	return c.MappingService.ToResponse(conf), nil
}

func (c *configControllerServiceImpl) Create(config *model.Config) (*response.Config, error) {
	conf := c.ConfigService.Set(config)
	if conf != nil {
		return nil, echo.NewHTTPError(400, fmt.Sprintf("%s already exists", config.Name))
	}

	return c.MappingService.ToResponse(conf), nil
}

func (c *configControllerServiceImpl) Update(config *model.Config) (*response.Config, error) {
	conf := c.ConfigService.Set(config)
	return c.MappingService.ToResponse(conf), nil
}

func (c *configControllerServiceImpl) Query(name string, query string) (*response.Config, error) {
	config, err := c.ConfigService.Query(name, query)
	if err != nil {
		return nil, err
	}

	return c.MappingService.ToResponse(config), nil
}
