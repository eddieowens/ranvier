package domain

import (
	"github.com/eddieowens/ranvier/commons"
)

type Schema struct {
	Name       string      `json:"name"`
	Extends    []string    `json:"extends" validate:"dive,ext,filepath"`
	Config     commons.Raw `json:"config" validate:"required_without=Name"`
	Type       string      `json:"type" validate:"oneof=yaml yml toml json"`
	Path       string      `json:"-"`
	IsAbstract bool        `json:"-"`
}
