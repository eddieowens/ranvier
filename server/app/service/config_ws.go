package service

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/labstack/gommon/log"
)

const ConfigWsServiceKey = "ConfigWsService"

type ConfigWsService interface {
	OnUpdate(eventType model.EventType, filepath string)
	OnStart(filepath string)
}

type configWsServiceImpl struct {
	ConfigService ConfigService `inject:"ConfigService"`
}

func (c *configWsServiceImpl) OnUpdate(eventType model.EventType, filepath string) {
	err := c.ConfigService.UpdateFromFile(eventType, filepath)
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
