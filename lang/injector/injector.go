package injector

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/semantics"
	"github.com/eddieowens/ranvier/lang/services"
)

var Injector axon.Injector

func CreateInjector() axon.Injector {
	binder := axon.NewBinder(
		new(compiler.Package),
		new(semantics.Package),
		new(services.Package),
	)

	return axon.NewInjector(binder)
}
