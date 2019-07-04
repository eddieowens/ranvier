package service

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/pubsub"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/oliveagle/jsonpath"
	"strings"
)

const ConfigServiceKey = "ConfigService"

type ConfigService interface {
	SetFromFile(filepath string) error
	UpdateFromFile(filepath string) error
	Delete(name string) *model.Config
	Set(config *model.Config) *model.Config
	Query(name string, query string) (*model.Config, error)
}

type configServiceImpl struct {
	ConfigMap collections.ConfigMap `inject:"ConfigMap"`
	PubSub    pubsub.PubSub         `inject:"PubSub"`
	Compiler  compiler.Compiler     `inject:"Compiler"`
	Config    configuration.Config  `inject:"Config"`
}

func (c *configServiceImpl) Set(config *model.Config) *model.Config {
	name := strings.ToLower(config.Name)
	_ = c.ConfigMap.WithLockWindow(name, func(_ model.Config, _ bool, saver collections.Saver) error {
		saver(*config)
		c.PubSub.Publish(config.Name, config)
		return nil
	})
	return config
}

func (c *configServiceImpl) Delete(name string) *model.Config {
	var conf *model.Config
	_ = c.ConfigMap.WithLock(func(configs map[string]model.Config) error {
		cfg, exists := configs[name]
		if exists {
			delete(configs, name)
			conf = &cfg
			c.PubSub.Publish(strings.ToLower(name), conf)
		}
		return nil
	})
	return conf
}

func (c *configServiceImpl) UpdateFromFile(filepath string) error {
	s, config, err := c.setFromFile(filepath)
	if err != nil {
		return err
	}

	if s.IsAbstract {
		return nil
	}

	c.PubSub.Publish(s.Name, config)

	return nil
}

func (c *configServiceImpl) SetFromFile(filepath string) error {
	_, _, err := c.setFromFile(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (c *configServiceImpl) setFromFile(filepath string) (*domain.Schema, *model.Config, error) {
	s, err := c.Compiler.Compile(filepath, &compiler.CompileOptions{
		OutputDirectory: c.Config.Compiler.OutputDirectory,
	})

	if err != nil {
		return nil, nil, err
	}

	if s.IsAbstract {
		return nil, nil, nil
	}

	var data interface{}
	_ = json.Unmarshal(s.Config, &data)

	config := model.Config{
		Name: s.Name,
		Data: data,
	}

	c.ConfigMap.Set(strings.ToLower(config.Name), config)

	return s, &config, nil
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
