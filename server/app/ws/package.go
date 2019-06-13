package ws

import (
	"github.com/eddieowens/axon"
)

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(WebsocketKey).To().Factory(websocketFactory).WithoutArgs(),
	}
}
