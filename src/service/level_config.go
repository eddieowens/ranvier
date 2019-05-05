package service

import (
	"config-manager/src/collections"
	"config-manager/src/model"
	"config-manager/src/state"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"net/http"
)

const LevelConfigServiceKey = "LevelConfigService"

type LevelConfigService interface {
	Update(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error)
	Query(level model.Level, id model.Id, query string) (*model.LevelConfig, error)
	Create(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error)
	Rollback(level model.Level, id model.Id, version int) (config model.LevelConfig, err error)
	Exists(level model.Level, id model.Id) bool
}

type levelConfigServiceImpl struct {
	State        state.LevelConfigState `inject:"LevelConfigState"`
	FileService  FileService            `inject:"FileService"`
	LevelService state.LevelService     `inject:"LevelService"`
	IdService    state.IdService        `inject:"IdService"`
	Json         jsoniter.API           `inject:"Json"`
}

func (l *levelConfigServiceImpl) Update(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error) {
	err = l.State.WithLock(level, id, func(levelConfig model.LevelConfig, exists bool, _ state.Saver) error {
		if exists {
			newConfig := make(map[string]interface{})
			oldConfig := make(map[string]interface{})
			err = l.Json.Unmarshal(data, &newConfig)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid json: %s", err.Error()))
			}

			_ = l.Json.Unmarshal(levelConfig.Config.GetRaw(), &oldConfig)

			err = mergo.Merge(&oldConfig, newConfig, mergo.WithOverride)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to update: %s", err.Error()))
			}

			mergedData, _ := l.Json.Marshal(oldConfig)

			config = l.newLevelConfig(level, id, mergedData)
			config.Version = levelConfig.Version + 1
		} else {
			config = l.newLevelConfig(level, id, data)
		}

		return l.create(config)
	})

	return config, err
}

func (l *levelConfigServiceImpl) Exists(level model.Level, id model.Id) bool {
	_, exists := l.State.Get(level, id)
	return exists
}

func (l *levelConfigServiceImpl) Rollback(level model.Level, id model.Id, version int) (config model.LevelConfig, err error) {
	versioned, exists := l.State.Get(level, l.IdService.VersionedId(id, version))
	if !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s version %d does not exist", l.IdService.Name(id), version))
	}
	err = l.State.WithLock(level, id, func(levelConfig model.LevelConfig, _ bool, _ state.Saver) error {
		if levelConfig.Version == version {
			config = levelConfig
			return nil
		}
		levelConfig.Config = versioned.Config
		levelConfig.Version = levelConfig.Version + 1
		config = levelConfig

		return l.FileService.Create(config)
	})

	return
}

func (l *levelConfigServiceImpl) Create(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error) {
	_, exists := l.State.Get(level, id)

	if !exists {
		err := l.State.WithLock(level, id, func(levelConfig model.LevelConfig, exists bool, _ state.Saver) error {
			if exists {
				return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("%s already exists", id))
			}

			config = l.newLevelConfig(level, id, data)

			return l.create(config)
		})
		if err != nil {
			return config, err
		}
	} else {
		err = echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("%s already exists", id))
	}

	return
}

func (l *levelConfigServiceImpl) Query(level model.Level, id model.Id, query string) (*model.LevelConfig, error) {
	var data model.LevelConfig
	if query == "" {
		data, _ = l.State.Get(level, id)
	} else {
		var ok bool
		if data, ok = l.State.Query(level, id, query); !ok {
			return nil, echo.NewHTTPError(
				http.StatusNotFound,
				fmt.Sprintf("the %s key could not be found", query),
			)
		}
	}
	return &data, nil
}

func (l *levelConfigServiceImpl) newLevelConfig(level model.Level, id model.Id, data []byte) model.LevelConfig {
	config := model.LevelConfig{
		Id:      id,
		Name:    l.IdService.Name(id),
		Version: 1,
		Level:   level,
		Config:  collections.NewJsonMap(data),
	}

	return config
}

func (l *levelConfigServiceImpl) create(levelConfig model.LevelConfig) error {
	if err := l.FileService.Create(levelConfig); err != nil {
		if sErr, ok := err.(*FileExistsError); !ok {
			return echo.NewHTTPError(http.StatusNotFound, sErr.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
