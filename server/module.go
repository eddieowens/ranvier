package server

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(AppKey).To().Instance(axon.StructPtr(new(appImpl))),
	}
}
