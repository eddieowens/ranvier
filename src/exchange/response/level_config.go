package response

type LevelConfig struct {
	Data *LevelConfigData `json:"data"`
}

type LevelConfigData struct {
	Version int         `json:"version"`
	Name    string      `json:"name"`
	Config  interface{} `json:"config"`
}

type LevelConfigMeta struct {
	Data *LevelConfigMetaData `json:"data"`
}

type LevelConfigMetaData struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type ApplicationsLevelConfigMeta struct {
	Data *ApplicationsLevelConfigMetaData `json:"data"`
}

type ApplicationsLevelConfigMetaData struct {
	Global       *LevelConfigMetaData  `json:"global"`
	Cluster      *LevelConfigMetaData  `json:"cluster"`
	Namespace    *LevelConfigMetaData  `json:"namespace"`
	Applications []LevelConfigMetaData `json:"applications"`
}

type ApplicationLevelConfigMeta struct {
	Data *ApplicationLevelConfigMetaData `json:"data"`
}

type ApplicationLevelConfigMetaData struct {
	Global      *LevelConfigMetaData `json:"global"`
	Cluster     *LevelConfigMetaData `json:"cluster"`
	Namespace   *LevelConfigMetaData `json:"namespace"`
	Application *LevelConfigMetaData `json:"application"`
}

type NamespacesLevelConfigMeta struct {
	Data *NamespacesLevelConfigMetaData `json:"data"`
}

type NamespacesLevelConfigMetaData struct {
	Global     *LevelConfigMetaData  `json:"global"`
	Cluster    *LevelConfigMetaData  `json:"cluster"`
	Namespaces []LevelConfigMetaData `json:"namespaces"`
}

type NamespaceLevelConfigMeta struct {
	Data *NamespaceLevelConfigMetaData `json:"data"`
}

type NamespaceLevelConfigMetaData struct {
	Global    *LevelConfigMetaData `json:"global"`
	Cluster   *LevelConfigMetaData `json:"cluster"`
	Namespace *LevelConfigMetaData `json:"namespace"`
}

type ClustersLevelConfigMeta struct {
	Data *ClustersLevelConfigMetaData `json:"data"`
}

type ClustersLevelConfigMetaData struct {
	Global   *LevelConfigMetaData  `json:"global"`
	Clusters []LevelConfigMetaData `json:"clusters"`
}

type ClusterLevelConfigMeta struct {
	Data *ClusterLevelConfigMetaData `json:"data"`
}

type ClusterLevelConfigMetaData struct {
	Global  *LevelConfigMetaData `json:"global"`
	Cluster *LevelConfigMetaData `json:"cluster"`
}
