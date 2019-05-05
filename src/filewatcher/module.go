package filewatcher

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(fileWatcherKey).To().Factory(fileWatcherFactory).WithoutArgs(),
	}
}
