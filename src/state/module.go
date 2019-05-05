package state

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(LevelConfigStateKey).To().Factory(levelConfigStateFactory).WithoutArgs(),
		axon.Bind(IdServiceKey).To().Instance(axon.StructPtr(new(idServiceImpl))),
		axon.Bind(LevelServiceKey).To().Instance(axon.StructPtr(new(levelServiceImpl))),
		axon.Bind(GlobalStateKey).To().Factory(levelConfigMapFactory).WithoutArgs(),
		axon.Bind(ClusterStateKey).To().Factory(levelConfigMapFactory).WithoutArgs(),
		axon.Bind(NamespaceStateKey).To().Factory(levelConfigMapFactory).WithoutArgs(),
		axon.Bind(ApplicationStateKey).To().Factory(levelConfigMapFactory).WithoutArgs(),
	}
}
