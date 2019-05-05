package service

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(FileServiceKey).To().Instance(axon.StructPtr(new(fileServiceImpl))),
		axon.Bind(LevelConfigServiceKey).To().Instance(axon.StructPtr(new(levelConfigServiceImpl))),
		axon.Bind(GlobalLevelConfigServiceKey).To().Instance(axon.StructPtr(new(globalLevelConfigServiceImpl))),
		axon.Bind(ClusterLevelConfigServiceKey).To().Instance(axon.StructPtr(new(clusterLevelConfigServiceImpl))),
		axon.Bind(NamespaceLevelConfigServiceKey).To().Instance(axon.StructPtr(new(namespaceLevelConfigServiceImpl))),
		axon.Bind(ApplicationLevelConfigServiceKey).To().Instance(axon.StructPtr(new(applicationLevelConfigServiceImpl))),
	}
}
