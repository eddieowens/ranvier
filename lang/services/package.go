package services

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons"
)

const JsonMergerKey = "JsonMerger"

const FilerKey = "Filer"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(FileCollectorKey).To().Instance(axon.StructPtr(new(fileCollectorImpl))),
		axon.Bind(FileFilterKey).To().Instance(axon.StructPtr(new(fileFilterImpl))),
		axon.Bind(FileServiceKey).To().Instance(axon.StructPtr(new(fileServiceImpl))),
		axon.Bind(JsonMergerKey).To().Instance(axon.StructPtr(commons.NewJsonMerger())),
		axon.Bind(FilerKey).To().Instance(axon.StructPtr(commons.NewFiler())),
		axon.Bind(ValidatorKey).To().Factory(ValidatorFactory).WithoutArgs(),
	}
}
