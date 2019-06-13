package poller

import "github.com/eddieowens/axon"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(GitPollerKey).To().Instance(axon.StructPtr(new(gitPollerImpl))),
	}
}
