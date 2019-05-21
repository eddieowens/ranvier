package model

import (
	"github.com/two-rabbits/ranvier/src/collections"
)

type Config struct {
	Name   string              `json:"name"`
	Config collections.JsonMap `json:"config"`
}

func (l Config) Copy() interface{} {
	return Config{
		Config: l.Config,
		Name:   l.Name,
	}
}
