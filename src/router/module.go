package router

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(RouterKey).To().Instance(axon.StructPtr(new(routerImpl))),
	}
}
