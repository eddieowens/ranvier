package app

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/configuration"
	"github.com/eddieowens/ranvier/server/app/controller"
	"github.com/eddieowens/ranvier/server/app/poller"
	"github.com/eddieowens/ranvier/server/app/pubsub"
	"github.com/eddieowens/ranvier/server/app/router"
	"github.com/eddieowens/ranvier/server/app/service"
	"github.com/eddieowens/ranvier/server/app/state"
	"github.com/eddieowens/ranvier/server/app/ws"
)

var Injector axon.Injector

func CreateInjector() axon.Injector {

	binder := axon.NewBinder(
		new(Package),
		new(controller.Package),
		new(service.Package),
		new(state.Package),
		new(poller.Package),
		new(router.Package),
		new(configuration.Package),
		new(ws.Package),
		new(pubsub.Package),
	)

	return axon.NewInjector(binder)
}
