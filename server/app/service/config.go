package service

import (
	"fmt"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/eddieowens/ranvier/server/app/pubsub"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/oliveagle/jsonpath"
	"path/filepath"
	"strings"
)

const ConfigServiceKey = "ConfigService"

type ConfigService interface {
	SetFromFile(filepath string) error
	UpdateFromFile(eventType model.EventType, filepath string) error
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
	_ = c.ConfigMap.WithLockWindow(name, func(_ model.Config, exists bool, saver collections.Saver) error {
		saver(name, *config)
		eventType := model.EventTypeCreate
		if exists {
			eventType = model.EventTypeUpdate
		}
		c.PubSub.Publish(config.Name, &model.ConfigEvent{
			EventType: eventType,
			Config:    *config,
		})
		return nil
	})
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

func (c *configServiceImpl) UpdateFromFile(eventType model.EventType, fp string) error {
	switch eventType {
	case model.EventTypeCreate, model.EventTypeUpdate:
		config, err := c.setFromFile(fp)
		if err != nil {
			return err
		}

		c.PubSub.Publish(config.Name, &model.ConfigEvent{
			EventType: eventType,
			Config:    *config,
		})
	case model.EventTypeDelete:
		c.Delete(compiler.ToSchemaName(fp))
	}

	return nil
}

func (c *configServiceImpl) SetFromFile(filepath string) error {
	_, err := c.setFromFile(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (c *configServiceImpl) setFromFile(fp string) (conf *model.Config, err error) {
	fp, _ = filepath.Rel(c.Config.Git.Directory, fp)
	s, err := c.Compiler.Compile(fp, compiler.CompileOptions{
		ParseOptions: compiler.ParseOptions{
			Root: c.Config.Git.Directory,
		},
		OutputDirectory: c.Config.Compiler.OutputDirectory,
	})

	if err != nil {
		return nil, err
	}

	var data interface{}
	_ = json.Unmarshal(s.Config, &data)

	config := model.Config{
		Name: s.Name,
		Data: data,
	}

	c.ConfigMap.Set(strings.ToLower(config.Name), config)

	return &config, nil
}

func (c *configServiceImpl) deleteFromFile(fp string) *model.Config {
	fp, _ = filepath.Rel(c.Config.Git.Directory, fp)
	var conf model.Config
	name := strings.ToLower(compiler.ToSchemaName(fp))
	_ = c.ConfigMap.WithLock(func(configs map[string]model.Config) error {
		conf = configs[name]
		delete(configs, name)
		return nil
	})
	return &conf
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
