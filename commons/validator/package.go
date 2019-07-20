package validator

import "github.com/eddieowens/axon"

const Key = "Validator"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(Key).To().Factory(Factory).WithoutArgs(),
	}
}
