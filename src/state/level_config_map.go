package state

import (
	"config-manager/src/collections"
	"config-manager/src/model"
	"github.com/eddieowens/axon"
	"github.com/json-iterator/go"
	"sync"
)

type Runner func(config model.LevelConfig, exists bool, saver Saver) error

type Saver func(config model.LevelConfig)

type LevelConfigMap interface {
	Query(id model.Id, query string) (config model.LevelConfig, exists bool)
	Set(levelConfig model.LevelConfig)
	Get(id model.Id) (model.LevelConfig, bool)
	WithLock(id model.Id, runner Runner) error
}

func levelConfigMapFactory(_ axon.Args) axon.Instance {
	return axon.StructPtr(&levelConfigMapImpl{
		m:    make(map[model.Id]model.LevelConfig),
		lock: sync.RWMutex{},
	})
}

type levelConfigMapImpl struct {
	m    map[model.Id]model.LevelConfig
	lock sync.RWMutex
	Json jsoniter.API `inject:"Json"`
}

func (s *levelConfigMapImpl) WithLock(id model.Id, runner Runner) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[id]
	return runner(val, exists, func(config model.LevelConfig) {
		s.m[config.Id] = config
	})
}

func (s *levelConfigMapImpl) Get(id model.Id) (model.LevelConfig, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	val, exists := s.m[id]
	return val, exists
}

func (s *levelConfigMapImpl) Set(levelConfig model.LevelConfig) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[levelConfig.Id] = levelConfig
}

func (s *levelConfigMapImpl) Query(id model.Id, query string) (config model.LevelConfig, exists bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if config, exists = s.m[id]; exists {
		raw, ok := config.Config.Get(query)
		if !ok {
			return config, ok
		}
		config = config.Copy().(model.LevelConfig)
		b, _ := s.Json.Marshal(raw)
		config.Config = collections.NewJsonMap(b)
		return
	} else {
		return config, false
	}
}
