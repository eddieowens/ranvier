package response

type LevelConfigResponse struct {
	Data LevelConfigData
}

type LevelConfigData struct {
	Version string      `json:"version"`
	Config  interface{} `json:"config"`
}

type LevelConfigMeta struct {
	Data LevelConfigMetaData `json:"data"`
}

type LevelConfigMetaData struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type ApplicationLevelConfigMeta struct {
	Data ApplicationLevelConfigMetaData `json:"data"`
}

type ApplicationLevelConfigMetaData struct {
	Global       LevelConfigMetaData   `json:"global"`
	Cluster      LevelConfigMetaData   `json:"cluster"`
	Namespace    LevelConfigMetaData   `json:"namespace"`
	Applications []LevelConfigMetaData `json:"applications"`
}

type NamespaceLevelConfigMeta struct {
	Data NamespaceLevelConfigMetaData `json:"data"`
}

type NamespaceLevelConfigMetaData struct {
	Global     LevelConfigMetaData   `json:"global"`
	Cluster    LevelConfigMetaData   `json:"cluster"`
	Namespaces []LevelConfigMetaData `json:"namespaces"`
}

type ClusterLevelConfigMeta struct {
	Data ClusterLevelConfigMetaData `json:"data"`
}

type ClusterLevelConfigMetaData struct {
	Global   LevelConfigMetaData   `json:"global"`
	Clusters []LevelConfigMetaData `json:"clusters"`
}

type GlobalLevelConfigMeta struct {
	Data LevelConfigMetaData `json:"data"`
}

type GlobalLevelConfigMetaData struct {
	Global LevelConfigMetaData `json:"global"`
}
