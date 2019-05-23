package configuration

import (
	"github.com/eddieowens/axon"
	"github.com/json-iterator/go"
)

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigKey).To().Factory(configFactory).WithoutArgs(),
		axon.Bind(JsonKey).To().Instance(axon.StructPtr(jsoniter.ConfigCompatibleWithStandardLibrary)),
	}
}
