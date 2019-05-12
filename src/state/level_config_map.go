package state

import (
	"github.com/eddieowens/axon"
	"github.com/json-iterator/go"
	"github.com/two-rabbits/ranvier/src/model"
	"sync"
)

type WriteRunnerWindow func(config model.LevelConfig, exists bool, saver Saver) error
type WriteRunner func(map[model.Id]model.LevelConfig) error
type ReadRunnerWindow func(config model.LevelConfig, exists bool) error
type Saver func(config model.LevelConfig)
type GetAllFilter func(id model.Id, config model.LevelConfig) bool

type LevelConfigMap interface {
	Set(levelConfig model.LevelConfig)
	Get(id model.Id) (model.LevelConfig, bool)
	GetAll(filter GetAllFilter) []model.LevelConfig
	Delete(id model.Id)
	WithLock(runner WriteRunner) error
	WithLockWindow(id model.Id, runner WriteRunnerWindow) error
	WithReadLockWindow(id model.Id, runner ReadRunnerWindow) error
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

func (s *levelConfigMapImpl) WithLock(runner WriteRunner) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return runner(s.m)
}

func (s *levelConfigMapImpl) Delete(id model.Id) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, id)
}

func (s *levelConfigMapImpl) GetAll(filter GetAllFilter) []model.LevelConfig {
	s.lock.RLock()
	defer s.lock.RUnlock()
	sl := make([]model.LevelConfig, 0)
	for k, v := range s.m {
		if filter == nil {
			sl = append(sl, v)
		} else if filter(k, v) {
			sl = append(sl, v)
		}
	}
	return sl
}

func (s *levelConfigMapImpl) WithLockWindow(id model.Id, runner WriteRunnerWindow) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[id]
	return runner(val, exists, func(config model.LevelConfig) {
		s.m[config.Id] = config
	})
}

func (s *levelConfigMapImpl) WithReadLockWindow(id model.Id, runner ReadRunnerWindow) error {
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
