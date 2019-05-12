package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/two-rabbits/ranvier/src/collections"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/state"
	"net/http"
)

const LevelConfigServiceKey = "LevelConfigService"

type LevelConfigService interface {
	Update(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error)
	Query(level model.Level, id model.Id, query string) (config model.LevelConfig, err error)
	MergedQuery(level model.Level, id model.Id, query string) (config model.LevelConfig, err error)
	Create(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error)
	Rollback(level model.Level, id model.Id, version int) (config model.LevelConfig, err error)
	Exists(level model.Level, id model.Id) bool
	GetAll(level model.Level) []model.LevelConfig
	Delete(level model.Level, id model.Id) (model.LevelConfig, error)
}

type levelConfigServiceImpl struct {
	State                   state.LevelConfigState        `inject:"LevelConfigState"`
	FileService             FileService                   `inject:"FileService"`
	LevelService            state.LevelService            `inject:"LevelService"`
	IdService               state.IdService               `inject:"IdService"`
	LevelConfigQueryService state.LevelConfigQueryService `inject:"LevelConfigQueryService"`
	MergeService            MergeService                  `inject:"MergeService"`
	Json                    jsoniter.API                  `inject:"Json"`
}

func (l *levelConfigServiceImpl) Delete(level model.Level, id model.Id) (model.LevelConfig, error) {
	resp, exists := l.State.Get(level, id)
	og := l.IdService.IdNames(id)
	name := l.IdService.Name(id)
	if !exists {
		levelStr := l.LevelService.ToString(level)
		return resp, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no %s found named %s", levelStr, name))
	}
	if level == model.Cluster || level == model.Namespace {
		for i := level; i <= model.Application; i++ {
			err := l.State.WithLock(i, func(configs map[model.Id]model.LevelConfig) error {
				for k, v := range configs {
					idNames := l.IdService.IdNames(k)
					shouldDelete := false
					switch i {
					case model.Namespace:
						shouldDelete = idNames.Namespace != "" && idNames.Namespace == og.Namespace
					case model.Cluster:
						shouldDelete = idNames.Cluster != "" && idNames.Cluster == og.Cluster

					}
					shouldDelete = shouldDelete && !l.IdService.IsVersionedId(k)
					if shouldDelete {
						err := l.FileService.Delete(v)
						if err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
						}
					}
				}
				return nil
			})
			if err != nil {
				return resp, err
			}
		}
	}
	return resp, nil

}

func (l *levelConfigServiceImpl) GetAll(level model.Level) []model.LevelConfig {
	return l.State.GetAll(level)
}

func (l *levelConfigServiceImpl) MergedQuery(level model.Level, id model.Id, query string) (config model.LevelConfig, err error) {
	levelInt := int(level)
	names := l.IdService.Names(id)
	mergedConfig := collections.NewJsonMap([]byte("{}"))
	var levelConfig model.LevelConfig
	for i := 0; i <= levelInt; i++ {
		xLevel := model.Level(i)
		exists := false
		if i == 0 {
			levelConfig, exists = l.State.Get(xLevel, state.GlobalId)
		} else {
			levelConfig, exists = l.State.Get(xLevel, l.IdService.Id(names[:i]...))
		}

		if exists {
			mergedConfig = l.MergeService.MergeJsonMaps(&mergedConfig, &levelConfig.Config)
		}
	}
	config = levelConfig
	config.Config = mergedConfig

	if data, exists := l.LevelConfigQueryService.Query(config, query); !exists {
		return config, echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("the %s key could not be found", query),
		)
	} else {
		return data, nil
	}
}

func (l *levelConfigServiceImpl) Update(level model.Level, id model.Id, data []byte) (config model.LevelConfig, err error) {
	err = l.State.WithLockWindow(level, id, func(levelConfig model.LevelConfig, exists bool, _ state.Saver) error {
		if exists {
			newConfig := collections.NewJsonMap(data)
			mergedData := l.MergeService.MergeJsonMaps(&levelConfig.Config, &newConfig)

			config = l.newLevelConfig(level, id, mergedData.GetRaw())
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
	err = l.State.WithLockWindow(level, id, func(levelConfig model.LevelConfig, _ bool, _ state.Saver) error {
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
		err := l.State.WithLockWindow(level, id, func(levelConfig model.LevelConfig, exists bool, _ state.Saver) error {
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

func (l *levelConfigServiceImpl) Query(level model.Level, id model.Id, query string) (config model.LevelConfig, err error) {
	if config, ok := l.State.Query(level, id, query); !ok {
		return config, echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("the %s key could not be found", query),
		)
	} else {
		return config, nil
	}
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
