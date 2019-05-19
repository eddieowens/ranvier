package compiler

import "github.com/eddieowens/axon"

type ConfigModule struct {
}

func (*ConfigModule) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigCompilerKey).To().Instance(axon.StructPtr(new(configCompilerImpl))),
	}
}
