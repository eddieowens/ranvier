package service

import "github.com/labstack/gommon/log"

const ConfigWsServiceKey = "ConfigWsService"

type ConfigWsService interface {
	OnUpdate(filepath string)
	OnStart(filepath string)
}

type configWsServiceImpl struct {
	ConfigService ConfigService `inject:"ConfigService"`
}

func (c *configWsServiceImpl) OnUpdate(filepath string) {
	err := c.ConfigService.UpdateFromFile(filepath)
	if err != nil {
		log.Warn(err)
	}
}

func (c *configWsServiceImpl) OnStart(filepath string) {
	err := c.ConfigService.SetFromFile(filepath)
	if err != nil {
		log.Warn(err)
	}
}
