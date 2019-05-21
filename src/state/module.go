package state

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigQueryServiceKey).To().Instance(axon.StructPtr(new(configQueryServiceImpl))),
		axon.Bind(ConfigMapKey).To().Instance(axon.StructPtr(new(configMapImpl))),
	}
}
