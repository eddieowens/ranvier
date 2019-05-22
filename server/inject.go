package server

import (
	"github.com/eddieowens/axon"
	"github.com/two-rabbits/ranvier/server/configuration"
	"github.com/two-rabbits/ranvier/server/controller"
	"github.com/two-rabbits/ranvier/server/poller"
	"github.com/two-rabbits/ranvier/server/router"
	"github.com/two-rabbits/ranvier/server/service"
	"github.com/two-rabbits/ranvier/server/state"
)

var Injector axon.Injector

func CreateInjector() axon.Injector {

	binder := axon.NewBinder(
		new(Module),
		new(controller.Module),
		new(service.Module),
		new(state.Module),
		new(poller.Module),
		new(router.Module),
		new(configuration.Module),
	)

	return axon.NewInjector(binder)
}
