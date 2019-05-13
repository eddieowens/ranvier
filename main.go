package main

import (
	"github.com/two-rabbits/ranvier/src"
)

// @title Ranvier
// @version 0.1
// @description A dynamic application configuration management system
// @termsOfService http://swagger.io/terms/
// @contact.name Edward Owens
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
func main() {
	src.Injector = src.CreateInjector()

	app := src.Injector.Get(src.AppKey).GetStructPtr().(src.App)

	app.Run()
}
