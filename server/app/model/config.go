package model

import "github.com/eddieowens/ranvier/lang/domain"

type Config struct {
	Name       string        `json:"name"`
	Data       interface{}   `json:"data"`
	Schema     domain.Schema `json:"-"`
	Dependents []string      `json:"-"`
}

type EventType int

const (
	EventTypeUpdate EventType = iota
	EventTypeCreate
	EventTypeDelete
)

type ConfigEvent struct {
	EventType EventType
	Config    Config
}

func (l Config) Copy() interface{} {
	return Config{
		Data: l.Data,
		Name: l.Name,
	}
}
