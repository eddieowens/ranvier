package main

import (
	"github.com/eddieowens/ranvier/ranvier/app"
	"log"
)

func main() {
	log.Fatal(app.NewInjector().GetStructPtr(app.Key).(app.App).Start())
}
