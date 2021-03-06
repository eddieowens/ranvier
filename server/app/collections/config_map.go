package collections

import (
	"github.com/eddieowens/ranvier/server/app/model"
	"sync"
)

type WriteRunnerWindow func(config model.Config, exists bool, saver Saver) error
type WriteRunner func(map[string]model.Config) error
type ReadRunnerWindow func(config model.Config, exists bool) error
type Saver func(name string, config model.Config)
type GetAllFilter func(name string, config model.Config) bool

type ConfigMap interface {
	Set(name string, levelConfig model.Config)
	Get(name string) (model.Config, bool)
	GetAll(filter GetAllFilter) []model.Config
	Delete(name string)
	WithLock(runner WriteRunner) error
	WithLockWindow(name string, runner WriteRunnerWindow) error
	WithReadLockWindow(name string, runner ReadRunnerWindow) error
}

func NewConfigMap() ConfigMap {
	return &configMapImpl{
		m:    make(map[string]model.Config),
		lock: sync.RWMutex{},
	}
}

type configMapImpl struct {
	m    map[string]model.Config
	lock sync.RWMutex
}

func (s *configMapImpl) WithLock(runner WriteRunner) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return runner(s.m)
}

func (s *configMapImpl) Delete(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, name)
}

func (s *configMapImpl) GetAll(filter GetAllFilter) []model.Config {
	s.lock.RLock()
	defer s.lock.RUnlock()
	sl := make([]model.Config, 0)
	for k, v := range s.m {
		if filter == nil {
			sl = append(sl, v)
		} else if filter(k, v) {
			sl = append(sl, v)
		}
	}
	return sl
}

func (s *configMapImpl) WithLockWindow(name string, runner WriteRunnerWindow) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[name]
	return runner(val, exists, func(name string, config model.Config) {
		s.m[name] = config
	})
}

func (s *configMapImpl) WithReadLockWindow(name string, runner ReadRunnerWindow) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, exists := s.m[name]
	return runner(val, exists)
}

func (s *configMapImpl) Get(name string) (model.Config, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	val, exists := s.m[name]
	return val, exists
}

func (s *configMapImpl) Set(name string, levelConfig model.Config) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[name] = levelConfig
}
