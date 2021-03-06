package semantics

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(AnalyzerKey).To().Instance(axon.StructPtr(new(analyzer))),
	}
}
