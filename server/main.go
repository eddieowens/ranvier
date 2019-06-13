package main

import "github.com/eddieowens/ranvier/server/app"

// @title Ranvier
// @version 0.1
// @description A dynamic application configuration management system
// @termsOfService http://swagger.io/terms/
// @contact.name Edward Owens
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
func main() {
	app.Injector = app.CreateInjector()

	app.Injector.Get(app.AppKey).GetStructPtr().(app.App).Run()
}
