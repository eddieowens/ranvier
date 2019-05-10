package response

type LevelConfigResponse struct {
	Data LevelConfigData
}

type LevelConfigData struct {
	Version string      `json:"version"`
	Config  interface{} `json:"config"`
}

type LevelConfigMeta struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type ApplicationsLevelConfigMeta struct {
	Global       LevelConfigMeta   `json:"global"`
	Cluster      LevelConfigMeta   `json:"cluster"`
	Namespace    LevelConfigMeta   `json:"namespace"`
	Applications []LevelConfigMeta `json:"applications"`
}

type ApplicationLevelConfigMeta struct {
	Global      LevelConfigMeta `json:"global"`
	Cluster     LevelConfigMeta `json:"cluster"`
	Namespace   LevelConfigMeta `json:"namespace"`
	Application LevelConfigMeta `json:"application"`
}

type NamespacesLevelConfigMeta struct {
	Global     LevelConfigMeta   `json:"global"`
	Cluster    LevelConfigMeta   `json:"cluster"`
	Namespaces []LevelConfigMeta `json:"namespaces"`
}

type NamespaceLevelConfigMeta struct {
	Global    LevelConfigMeta `json:"global"`
	Cluster   LevelConfigMeta `json:"cluster"`
	Namespace LevelConfigMeta `json:"namespace"`
}

type ClustersLevelConfigMeta struct {
	Global   LevelConfigMeta   `json:"global"`
	Clusters []LevelConfigMeta `json:"clusters"`
}

type ClusterLevelConfigMeta struct {
	Global  LevelConfigMeta `json:"global"`
	Cluster LevelConfigMeta `json:"cluster"`
}
