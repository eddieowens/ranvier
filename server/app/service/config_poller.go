package service

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/labstack/gommon/log"
)

const ConfigPollerServiceKey = "ConfigPollerService"

type ConfigPollerService interface {
	OnUpdate(eventType model.EventType, filepath string)
	OnStart(filepath string)
}

type configPollerServiceImpl struct {
	ConfigService ConfigService           `inject:"ConfigService"`
	ConfigDepMap  map[string]model.Config `inject:"ConfigDepMap"`
}

func (c *configPollerServiceImpl) OnUpdate(eventType model.EventType, filepath string) {
	err := c.ConfigService.UpdateFromFile(eventType, filepath)
	if err != nil {
		log.Warn(err)
	}
}

func (c *configPollerServiceImpl) OnStart(filepath string) {
	err := c.ConfigService.SetFromFile(filepath)
	if err != nil {
		log.Warn(err)
	}
}
