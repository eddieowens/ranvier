package compiler

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(Key).To().Instance(axon.StructPtr(new(compilerImpl))),
		axon.Bind(SchemaPackerKey).To().Instance(axon.StructPtr(new(schemaPackerImpl))),
	}
}
