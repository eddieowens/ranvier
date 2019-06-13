package state

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/collections"
)

const ConfigMapKey = "ConfigMap"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigMapKey).To().Instance(axon.StructPtr(collections.NewConfigMap())),
	}
}
