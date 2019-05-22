package main

import (
	"github.com/two-rabbits/ranvier/server"
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
	server.Injector = server.CreateInjector()

	app := server.Injector.Get(server.AppKey).GetStructPtr().(server.App)

	app.Run()
}
