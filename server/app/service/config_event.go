package service

import (
	"github.com/eddieowens/ranvier/server/app/model"
	log "github.com/sirupsen/logrus"
)

const ConfigEventServiceKey = "ConfigEventService"

type ConfigEventService interface {
	OnEvent(eventType model.EventType, filepath string)
	OnStart(filepath string)
}

type configEventServiceImpl struct {
	ConfigService   ConfigMapService        `inject:"ConfigMapService"`
	ConfigDepMap    map[string]model.Config `inject:"ConfigDepMap"`
	CompilerService CompilerService         `inject:"CompilerService"`
}

func (c *configEventServiceImpl) OnEvent(eventType model.EventType, filepath string) {
	conf, scheme, err := c.CompilerService.Compile(filepath)
	if err != nil {
		log.WithError(err).WithField("filepath", filepath).Error("Failed to compile file")
		return
	}
	switch eventType {
	case model.EventTypeCreate:
		c.ConfigService.Create(conf)
	case model.EventTypeUpdate:
		c.ConfigService.Update(conf)
	case model.EventTypeDelete:
		c.ConfigService.Delete(conf.Name)
	}
}

func (c *configEventServiceImpl) OnStart(filepath string) {
	conf, scheme, err := c.CompilerService.Compile(filepath)
	if err != nil {
		log.WithError(err).WithField("filepath", filepath).Error("Failed to compile file")
		return
	}
	c.ConfigService.Set(conf)
}
