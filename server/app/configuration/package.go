package configuration

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigKey).To().Factory(configFactory).WithoutArgs(),
		axon.Bind(AuthMethodKey).To().Factory(authMethodFactory).WithoutArgs(),
	}
}
