package controller

import "github.com/eddieowens/axon"

type Module struct {
}

func (*Module) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(GlobalControllerKey).To().Instance(axon.StructPtr(new(globalControllerImpl))),
		axon.Bind(ClusterControllerKey).To().Instance(axon.StructPtr(new(clusterControllerImpl))),
		axon.Bind(NamespaceControllerKey).To().Instance(axon.StructPtr(new(namespaceControllerImpl))),
		axon.Bind(ApplicationControllerKey).To().Instance(axon.StructPtr(new(applicationControllerImpl))),
	}
}
