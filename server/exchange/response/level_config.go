package response

type Config struct {
	Data *LevelConfigData `json:"data"`
}

type LevelConfigData struct {
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}

type ConfigMeta struct {
	Data *LevelConfigMetaData `json:"data"`
}

type LevelConfigMetaData struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}
