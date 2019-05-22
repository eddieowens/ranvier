package state

import (
	"github.com/eddieowens/axon"
)

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigMapKey).To().Factory(configMapFactory).WithoutArgs(),
	}
}
