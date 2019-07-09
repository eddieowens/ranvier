package pubsub

import "github.com/eddieowens/axon"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(Key).To().Factory(pubSubFactory).WithoutArgs(),
	}
}
