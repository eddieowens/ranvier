package app

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/ranvier/app/beans"
	"github.com/eddieowens/ranvier/ranvier/app/cmd"
	"github.com/eddieowens/ranvier/ranvier/app/cmd/compile"
)

func NewInjector() axon.Injector {
	return axon.NewInjector(axon.NewBinder(
		new(beans.Package),
		new(compile.Package),
		new(cmd.Package),
		new(Package),
	))
}
