package service

import (
	"fmt"
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/pubsub"
	"github.com/labstack/echo"
	"github.com/oliveagle/jsonpath"
	"strings"
)

const ConfigServiceKey = "ConfigService"

type ConfigService interface {
	Delete(name string) *model.Config
	Set(config *model.Config) *model.Config
	Update(config *model.Config) *model.Config
	Query(name string, query string) (*model.Config, error)
}

type configServiceImpl struct {
	ConfigMap collections.ConfigMap `inject:"ConfigMap"`
	PubSub    pubsub.PubSub         `inject:"PubSub"`
}

func (c *configServiceImpl) Update(config *model.Config) *model.Config {

}

func (c *configServiceImpl) Set(config *model.Config) *model.Config {
	name := strings.ToLower(config.Name)
	c.ConfigMap.Set(name, *config)
	return config
}

func (c *configServiceImpl) Delete(name string) *model.Config {
	var conf *model.Config
	_ = c.ConfigMap.WithLock(func(configs map[string]model.Config) error {
		name = strings.ToLower(name)
		cfg, exists := configs[name]
		if exists {
			delete(configs, name)
			conf = &cfg
			c.PubSub.Publish(strings.ToLower(name), &model.ConfigEvent{
				EventType: model.EventTypeDelete,
				Config:    cfg,
			})
		}
		return nil
	})
	return conf
}

func (c *configServiceImpl) Query(name string, query string) (*model.Config, error) {
	if query == "" {
		query = "$"
	}
	strings.ToLower(name)
	config, exists := c.ConfigMap.Get(name)
	if !exists {
		return nil, echo.NewHTTPError(404, fmt.Sprintf("config with name %s could not be found", name))
	}

	raw, err := jsonpath.JsonPathLookup(config.Data, query)
	if err != nil {
		return nil, echo.NewHTTPError(400, err.Error())
	}

	return &model.Config{
		Name: name,
		Data: raw,
	}, nil
}
