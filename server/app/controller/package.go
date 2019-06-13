package controller

import "github.com/eddieowens/axon"

const ControllersKey = "Controllers"

type Package struct {
}

func (*Package) Bindings() []axon.Binding {
	return []axon.Binding{
		axon.Bind(ConfigControllerKey).To().StructPtr(new(configControllerImpl)),
		axon.Bind(WebsocketControllerKey).To().StructPtr(new(websocketControllerImpl)),
		axon.Bind(ControllersKey).To().Keys(ConfigControllerKey, WebsocketControllerKey),
	}
}
