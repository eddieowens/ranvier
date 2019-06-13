package cli

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(CliKey).To().StructPtr(new(cliImpl)),
	}
}
