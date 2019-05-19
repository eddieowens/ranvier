package controller

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(GlobalControllerKey).To().Instance(axon.StructPtr(new(globalControllerImpl))),
	}
}
