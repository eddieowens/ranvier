package app

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(Key).To().StructPtr(new(app)),
	}
}
