package cmd

import "github.com/eddieowens/axon"

const CmdsKey = "Cmds"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(CompileCmdKey).To().StructPtr(new(compileCmd)),
		axon.Bind(CmdsKey).To().Keys(CompileCmdKey),
	}
}
