package poller

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(GitPollerKey).To().Instance(axon.StructPtr(new(gitPollerImpl))),
	}
}
