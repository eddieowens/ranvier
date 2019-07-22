package cmd

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/ranvier/app/cmd/compile"
)

const CommandsKey = "Commands"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(CommandsKey).To().Keys(compile.CmdKey),
	}
}
