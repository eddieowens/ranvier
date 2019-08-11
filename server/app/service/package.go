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
		axon.Bind(ConfigControllerServiceKey).To().StructPtr(new(configControllerServiceImpl)),
		axon.Bind(JsonMergerKey).To().StructPtr(commons.NewJsonMerger()),
		axon.Bind(MappingServiceKey).To().StructPtr(new(mappingServiceImpl)),
		axon.Bind(ConfigMapServiceKey).To().StructPtr(new(configMapServiceImpl)),
		axon.Bind(ConfigServiceKey).To().StructPtr(new(configServiceImpl)),
		axon.Bind(GitServiceKey).To().StructPtr(new(gitServiceImpl)),
		axon.Bind(ConfigEventServiceKey).To().StructPtr(new(configEventServiceImpl)),
		axon.Bind(CompilerServiceKey).To().StructPtr(new(compilerServiceImpl)),
	}
}
