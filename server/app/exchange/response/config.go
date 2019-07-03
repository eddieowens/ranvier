package response

import "github.com/eddieowens/ranvier/server/app/model"

type Config struct {
	Data *model.Config `json:"data"`
}

type ConfigData struct {
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}

type ConfigMeta struct {
	Data *ConfigMetaData `json:"data"`
}

type ConfigMetaData struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}
