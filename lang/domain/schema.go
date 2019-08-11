package domain

import (
	"github.com/eddieowens/ranvier/commons"
)

type Schema struct {
	// The name of the config schema. This name is generated based upon the dir and will be ignored if specified within
	// the schema file. All schema file names must begin and end with an alphanumeric character and can only contain the
	// '-' special character.
	Name string `json:"name" validate:"required,dns_1123"`

	// The path relative to the root path of the config. The root path of the config is typically the root path of the git
	// repo. These fields cannot contain the '.' or '..' characters.
	Extends []string `json:"extends" validate:"dive,ext=json toml yaml yml,file"`

	// The user's own config. This can be any data the user wishes to use.
	Config commons.Raw `json:"config"`

	// When using the Kubernetes plugin, this value determines what namespace to place the resulting configmap into. If
	// no Kubernetes plugin is running, this field has no effect.
	Namespace string `json:"namespace"`

	// The fullfile path to the file from which this schema was created. Ignored if specified by user as it is only meant
	// to be used internally.
	Path string `json:"-"`
}

type CompiledSchema struct {
	ParsedSchema
}

type ParsedSchema struct {
	Schema
	Dependencies []ParsedSchema
}
