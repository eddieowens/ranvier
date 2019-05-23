package service

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigControllerServiceKey).To().Instance(axon.StructPtr(new(configControllerServiceImpl))),
		axon.Bind(MergeServiceKey).To().Instance(axon.StructPtr(new(mergeServiceImpl))),
		axon.Bind(MappingServiceKey).To().Instance(axon.StructPtr(new(mappingServiceImpl))),
		axon.Bind(ConfigQueryServiceKey).To().Instance(axon.StructPtr(new(configQueryServiceImpl))),
	}
}
