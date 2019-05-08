package model

import (
	"github.com/two-rabbits/ranvier/src/collections"
)

type LevelConfig struct {
	Version int                 `json:"version"`
	Name    string              `json:"name"`
	Id      Id                  `json:"id"`
	Level   Level               `json:"level"`
	Config  collections.JsonMap `json:"config"`
}

func (l LevelConfig) Copy() interface{} {
	return LevelConfig{
		Version: l.Version,
		Config:  l.Config,
		Level:   l.Level,
		Name:    l.Name,
		Id:      l.Id,
	}
}

type LevelConfigMeta struct {
	Versions []LevelConfig `json:"versions"`
}

func (l *LevelConfigMeta) Copy() interface{} {
	return LevelConfigMeta{
		Versions: l.Versions,
	}
}
