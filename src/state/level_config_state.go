package state

import (
	"github.com/eddieowens/axon"
	"github.com/two-rabbits/ranvier/src/model"
	"sort"
	"sync"
)

const LevelConfigStateKey = "LevelConfigState"
const GlobalStateKey = "GlobalState"
const ClusterStateKey = "ClusterState"
const NamespaceStateKey = "NamespaceState"
const ApplicationStateKey = "ApplicationState"

const GlobalId = model.Id("global")

type LevelConfigState interface {
	Query(level model.Level, id model.Id, query string) (config model.LevelConfig, exists bool)
	Set(levelConfig model.LevelConfig)
	Get(level model.Level, id model.Id) (levelConfig model.LevelConfig, exists bool)
	GetAll(level model.Level) []model.LevelConfig
	Delete(level model.Level, id model.Id) (model.LevelConfig, bool)
	WithLockWindow(level model.Level, id model.Id, runner WriteRunnerWindow) error
	WithLock(level model.Level, runner WriteRunner) error
}

type levelConfigStateImpl struct {
	GlobalState             LevelConfigMap          `inject:"GlobalState"`
	ClusterState            LevelConfigMap          `inject:"ClusterState"`
	NamespaceState          LevelConfigMap          `inject:"NamespaceState"`
	ApplicationState        LevelConfigMap          `inject:"ApplicationState"`
	IdService               IdService               `inject:"IdService"`
	LevelConfigQueryService LevelConfigQueryService `inject:"LevelConfigQueryService"`
	VersionMap              map[model.Id][]int
	VersionMapLock          sync.RWMutex
}

func (l *levelConfigStateImpl) WithLock(level model.Level, runner WriteRunner) error {
	switch level {
	case model.Global:
		return l.GlobalState.WithLock(runner)
	case model.Cluster:
		return l.ClusterState.WithLock(runner)
	case model.Namespace:
		return l.NamespaceState.WithLock(runner)
	case model.Application:
		return l.ApplicationState.WithLock(runner)
	default:
		return nil
	}
}

func (l *levelConfigStateImpl) Delete(level model.Level, id model.Id) (model.LevelConfig, bool) {
	switch level {
	case model.Global:
		return l.delete(id, l.GlobalState)
	case model.Cluster:
		return l.delete(id, l.ClusterState)
	case model.Namespace:
		return l.delete(id, l.NamespaceState)
	case model.Application:
		return l.delete(id, l.ApplicationState)
	default:
		return model.LevelConfig{}, false
	}
}

func (l *levelConfigStateImpl) GetAll(level model.Level) []model.LevelConfig {
	switch level {
	case model.Global:
		return l.GlobalState.GetAll(l.getAllFilter)
	case model.Cluster:
		return l.ClusterState.GetAll(l.getAllFilter)
	case model.Namespace:
		return l.NamespaceState.GetAll(l.getAllFilter)
	case model.Application:
		return l.ApplicationState.GetAll(l.getAllFilter)
	default:
		return []model.LevelConfig{}
	}
}

func (l *levelConfigStateImpl) WithLockWindow(level model.Level, id model.Id, runner WriteRunnerWindow) error {
	switch level {
	case model.Global:
		return l.GlobalState.WithLockWindow(id, runner)
	case model.Cluster:
		return l.ClusterState.WithLockWindow(id, runner)
	case model.Namespace:
		return l.NamespaceState.WithLockWindow(id, runner)
	case model.Application:
		return l.ApplicationState.WithLockWindow(id, runner)
	default:
		return nil
	}
}

func (l *levelConfigStateImpl) Get(level model.Level, id model.Id) (levelConfig model.LevelConfig, exists bool) {
	switch level {
	case model.Global:
		return l.GlobalState.Get(id)
	case model.Cluster:
		return l.ClusterState.Get(id)
	case model.Namespace:
		return l.NamespaceState.Get(id)
	case model.Application:
		return l.ApplicationState.Get(id)
	default:
		return
	}
}

func (l *levelConfigStateImpl) Set(levelConfig model.LevelConfig) {
	switch levelConfig.Level {
	case model.Global:
		l.set(levelConfig, l.GlobalState)
	case model.Cluster:
		l.set(levelConfig, l.ClusterState)
	case model.Namespace:
		l.set(levelConfig, l.NamespaceState)
	case model.Application:
		l.set(levelConfig, l.ApplicationState)
	}
	return
}

func (l *levelConfigStateImpl) Query(level model.Level, id model.Id, query string) (config model.LevelConfig, exists bool) {
	switch level {
	case model.Global:
		return l.query(l.GlobalState, id, query)
	case model.Cluster:
		return l.query(l.ClusterState, id, query)
	case model.Namespace:
		return l.query(l.NamespaceState, id, query)
	case model.Application:
		return l.query(l.ApplicationState, id, query)
	default:
		return config, false
	}
}

func (l *levelConfigStateImpl) set(levelConfig model.LevelConfig, configMap LevelConfigMap) {
	_ = configMap.WithLockWindow(levelConfig.Id, func(_ model.LevelConfig, _ bool, saver Saver) error {
		versionedId := l.IdService.VersionedId(levelConfig.Id, levelConfig.Version)

		versions := l.VersionMap[levelConfig.Id]
		if versions == nil {
			versions = []int{}
		}
		sort.Ints(versions)
		id := levelConfig.Id
		levelConfig.Id = versionedId
		saver(levelConfig)
		if len(versions) == 0 || versions[len(versions)-1] < levelConfig.Version {
			levelConfig.Id = id
			saver(levelConfig)
		}
		versions = append(versions, levelConfig.Version)
		l.VersionMap[id] = versions
		return nil
	})
}

func (l *levelConfigStateImpl) getAllFilter(id model.Id, _ model.LevelConfig) bool {
	return !l.IdService.IsVersionedId(id)
}

func (l *levelConfigStateImpl) query(state LevelConfigMap, id model.Id, query string) (config model.LevelConfig, exists bool) {
	_ = state.WithReadLockWindow(id, func(levelConfig model.LevelConfig, _ bool) error {
		config, exists = l.LevelConfigQueryService.Query(levelConfig, query)
		return nil
	})
	return config, exists
}

func (l *levelConfigStateImpl) delete(id model.Id, state LevelConfigMap) (config model.LevelConfig, exists bool) {
	_ = state.WithLock(func(configs map[model.Id]model.LevelConfig) error {
		config, exists = configs[id]
		if exists {
			delete(configs, id)
			versions := l.VersionMap[id]
			for _, v := range versions {
				delete(configs, l.IdService.VersionedId(id, v))
			}
			delete(l.VersionMap, id)
		}
		return nil
	})
	return
}

func levelConfigStateFactory(_ axon.Args) axon.Instance {
	return axon.StructPtr(&levelConfigStateImpl{
		VersionMap:     make(map[model.Id][]int),
		VersionMapLock: sync.RWMutex{},
	})
}
