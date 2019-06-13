package app

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/lang"
)

type Package struct {
}

const CompilerKey = "Compiler"

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(AppKey).To().Instance(axon.StructPtr(new(appImpl))),
		axon.Bind(CompilerKey).To().Instance(axon.StructPtr(lang.NewCompiler())),
	}
}
