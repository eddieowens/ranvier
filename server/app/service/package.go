package service

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons"
)

const JsonMergerKey = "JsonMerger"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigControllerServiceKey).To().Instance(axon.StructPtr(new(configControllerServiceImpl))),
		axon.Bind(JsonMergerKey).To().Instance(axon.StructPtr(commons.NewJsonMerger())),
		axon.Bind(MappingServiceKey).To().Instance(axon.StructPtr(new(mappingServiceImpl))),
		axon.Bind(ConfigServiceKey).To().Instance(axon.StructPtr(new(configServiceImpl))),
		axon.Bind(GitServiceKey).To().Instance(axon.StructPtr(new(gitServiceImpl))),
	}
}
