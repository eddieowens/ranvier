package compile

import "github.com/eddieowens/axon"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(CmdKey).To().StructPtr(new(compileCmd)),
	}
}
