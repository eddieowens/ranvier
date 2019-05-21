package poller

import "github.com/eddieowens/axon"

type Poller struct {
}

func (*Poller) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(GitPollerKey).To().Instance(axon.StructPtr(new(gitPollerImpl))),
	}
}
