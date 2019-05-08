package src

import (
	"github.com/eddieowens/axon"
	"github.com/two-rabbits/ranvier/src/configuration"
	"github.com/two-rabbits/ranvier/src/controller"
	"github.com/two-rabbits/ranvier/src/filewatcher"
	"github.com/two-rabbits/ranvier/src/router"
	"github.com/two-rabbits/ranvier/src/service"
	"github.com/two-rabbits/ranvier/src/state"
)

var Injector axon.Injector

func CreateInjector() axon.Injector {

	binder := axon.NewBinder(
		new(Module),
		new(controller.Module),
		new(service.Module),
		new(filewatcher.Module),
		new(state.Module),
		new(router.Module),
		new(configuration.Module),
	)

	return axon.NewInjector(binder)
}
