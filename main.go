package main

import (
	"github.com/two-rabbits/ranvier/src"
)

func main() {
	src.Injector = src.CreateInjector()

	app := src.Injector.Get(src.AppKey).GetStructPtr().(src.App)

	app.Run()
}
