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

	// The config file extension type that will be output after compilation. Valid file extensions are
	//   * yaml
	//   * yml
	//   * toml
	//   * json
	// Defaults to json.
	Type string `json:"type" validate:"oneof=yaml yml toml json"`

	// The fullfile path to the file from which this schema was created. Ignored if specified by user as it is only meant
	// to be used internally.
	Path string `json:"-"`
}
