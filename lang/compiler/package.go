package compiler

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(CompilerKey).To().Instance(axon.StructPtr(new(compilerImpl))),
		axon.Bind(PackerKey).To().Instance(axon.StructPtr(new(packerImpl))),
	}
}
