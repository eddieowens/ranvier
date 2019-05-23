package app

import (
	"github.com/eddieowens/axon"
	"github.com/two-rabbits/ranvier/server/app/configuration"
	"github.com/two-rabbits/ranvier/server/app/controller"
	"github.com/two-rabbits/ranvier/server/app/poller"
	"github.com/two-rabbits/ranvier/server/app/router"
	"github.com/two-rabbits/ranvier/server/app/service"
	"github.com/two-rabbits/ranvier/server/app/state"
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
