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

const ConfigServiceKey = "ConfigService"

type ConfigService interface {
	Update(id model.Id, data []byte) (config *model.Config, err error)
	Query(id model.Id, query string) (config *model.Config, exists bool)
	MergedQuery(id model.Id, query string) (config *model.Config, err error)
	Create(id model.Id, data []byte) (config *model.Config, err error)
	Rollback(id model.Id, version int) (config *model.Config, err error)
	Exists(id model.Id) bool
	GetAll() []model.Config
	Delete(id model.Id) (*model.Config, error)
}

type levelConfigServiceImpl struct {
	State                   state.ConfigMap          `inject:"ConfigState"`
	FileService             FileService              `inject:"FileService"`
	LevelService            state.LevelService       `inject:"LevelService"`
	IdService               state.IdService          `inject:"IdService"`
	LevelConfigQueryService state.ConfigQueryService `inject:"ConfigQueryService"`
	MergeService            MergeService             `inject:"MergeService"`
	Json                    jsoniter.API             `inject:"Json"`
}

func (l *levelConfigServiceImpl) Delete(id model.Id) (*model.Config, error) {
	resp, exists := l.State.Get(level, id)
	og := l.IdService.IdNames(id)
	name := l.IdService.Name(id)
	if !exists {
		levelStr := l.LevelService.ToString(level)
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no %s found named %s", levelStr, name))
	}
	if level != model.Global && level != model.Application {
		for i := level; i <= model.Application; i++ {
			err := l.State.WithLock(i, func(configs map[model.Id]model.Config) error {
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
				return nil, err
			}
		}
	}

	err := l.FileService.Delete(resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (l *levelConfigServiceImpl) GetAll() []model.Config {
	return l.State.GetAll(level)
}

func (l *levelConfigServiceImpl) MergedQuery(id model.Id, query string) (config *model.Config, err error) {
	levelInt := int(level)
	names := l.IdService.Names(id)
	mergedConfig := collections.NewJsonMap([]byte("{}"))
	var levelConfig model.Config
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
	config = &levelConfig
	config.Config = mergedConfig

	if data, exists := l.LevelConfigQueryService.Query(*config, query); !exists {
		return config, echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("the %s key could not be found", query),
		)
	} else {
		return &data, nil
	}
}

func (l *levelConfigServiceImpl) Update(id model.Id, data []byte) (config *model.Config, err error) {
	err = l.State.WithLockWindow(level, id, func(levelConfig model.Config, exists bool, _ state.Saver) error {
		if exists {
			newConfig := collections.NewJsonMap(data)
			mergedData := l.MergeService.MergeJsonMaps(&levelConfig.Config, &newConfig)

			c := l.newLevelConfig(level, id, mergedData.GetRaw())
			config = &c
			config.Version = levelConfig.Version + 1
		} else {
			c := l.newLevelConfig(level, id, data)
			config = &c
		}

		return l.create(config)
	})

	return config, err
}

func (l *levelConfigServiceImpl) Exists(id model.Id) bool {
	_, exists := l.State.Get(level, id)
	return exists
}

func (l *levelConfigServiceImpl) Rollback(id model.Id, version int) (config *model.Config, err error) {
	versioned, exists := l.State.Get(level, l.IdService.VersionedId(id, version))
	if !exists {
		return config, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s version %d does not exist", l.IdService.Name(id), version))
	}
	err = l.State.WithLockWindow(level, id, func(levelConfig model.Config, _ bool, _ state.Saver) error {
		if levelConfig.Version == version {
			config = &levelConfig
			return nil
		}
		levelConfig.Config = versioned.Config
		levelConfig.Version = levelConfig.Version + 1
		config = &levelConfig

		return l.FileService.Create(*config)
	})

	return
}

func (l *levelConfigServiceImpl) Create(id model.Id, data []byte) (config *model.Config, err error) {
	_, exists := l.State.Get(level, id)

	if !exists {
		err := l.State.WithLockWindow(level, id, func(levelConfig model.Config, exists bool, _ state.Saver) error {
			if exists {
				return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("%s already exists", id))
			}

			c := l.newLevelConfig(level, id, data)
			config = &c

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

func (l *levelConfigServiceImpl) Query(id model.Id, query string) (config *model.Config, exists bool) {
	c, ok := l.State.Query(level, id, query)
	if !ok {
		return nil, ok
	}
	return &c, ok
}

func (l *levelConfigServiceImpl) newLevelConfig(level model.Level, id model.Id, data []byte) model.Config {
	config := model.Config{
		Id:      id,
		Name:    l.IdService.Name(id),
		Version: 1,
		Level:   level,
		Config:  collections.NewJsonMap(data),
	}

	return config
}

func (l *levelConfigServiceImpl) create(levelConfig *model.Config) error {
	if err := l.FileService.Create(*levelConfig); err != nil {
		if sErr, ok := err.(*FileExistsError); !ok {
			return echo.NewHTTPError(http.StatusNotFound, sErr.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
