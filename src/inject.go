package src

import (
	"config-manager/src/configuration"
	"config-manager/src/controller"
	"config-manager/src/filewatcher"
	"config-manager/src/router"
	"config-manager/src/service"
	"config-manager/src/state"
	"github.com/eddieowens/axon"
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
