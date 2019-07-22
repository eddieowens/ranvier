package beans

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/commons/validator"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(validator.Key).To().Factory(validator.Factory).WithoutArgs(),
	}
}
