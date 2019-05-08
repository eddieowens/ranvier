package state

import (
	"github.com/eddieowens/axon"
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/src/model"
	"sync"
)

type WriteRunner func(config model.LevelConfig, exists bool, saver Saver) error
type ReadRunner func(config model.LevelConfig, exists bool) error

type Saver func(config model.LevelConfig)

type LevelConfigMap interface {
	Set(levelConfig model.LevelConfig)
	Get(id model.Id) (model.LevelConfig, bool)
	WithLock(id model.Id, runner WriteRunner) error
	WithReadLock(id model.Id, runner ReadRunner) error
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

func (s *levelConfigMapImpl) WithLock(id model.Id, runner WriteRunner) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[id]
	return runner(val, exists, func(config model.LevelConfig) {
		s.m[config.Id] = config
	})
}

func (s *levelConfigMapImpl) WithReadLock(id model.Id, runner ReadRunner) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[id]
	return runner(val, exists)
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
