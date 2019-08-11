package service

import (
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/model"
	"strings"
)

const ConfigMapServiceKey = "ConfigMapService"

type ConfigMapService interface {
	Get(name string) *model.Config
	Delete(name string)
	Set(config *model.Config)
}

type configMapServiceImpl struct {
	ConfigMap collections.ConfigMap `inject:"ConfigMap"`
}

func (c *configMapServiceImpl) Get(name string) *model.Config {
	conf, exists := c.ConfigMap.Get(strings.ToLower(name))
	if exists {
		return &conf
	}
	return nil
}

func (c *configMapServiceImpl) Delete(name string) {
	c.ConfigMap.Delete(strings.ToLower(name))
}

func (c *configMapServiceImpl) Set(config *model.Config) {
	c.ConfigMap.Set(strings.ToLower(config.Name), *config)
}
