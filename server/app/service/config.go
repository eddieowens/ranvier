package service

import (
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/model"
	"strings"
)

const ConfigServiceKey = "ConfigService"

type ConfigService interface {
	Set(config *model.Config) error
	Get(config *model.Config) *model.Config
	Delete(name string) error
	Query(name, query string) (*model.Config, error)
}

type configServiceImpl struct {
	ConfigMap collections.ConfigMap `inject:"ConfigMap"`
}

func (c *configServiceImpl) Set(config *model.Config) error {
	return c.ConfigMap.WithLock(func(configs map[string]model.Config) error {
		og, exists := configs[strings.ToLower(config.Name)]
		if exists {

		} else {

		}
		return nil
	})
}

func (c *configServiceImpl) Get(config *model.Config) *model.Config {
	panic("implement me")
}

func (c *configServiceImpl) Delete(name string) error {
	panic("implement me")
}

func (c *configServiceImpl) Query(name, query string) (*model.Config, error) {
	panic("implement me")
}
