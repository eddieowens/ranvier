package router

import "github.com/eddieowens/axon"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(RouterKey).To().Instance(axon.StructPtr(new(routerImpl))),
	}
}
