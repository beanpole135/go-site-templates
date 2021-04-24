package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

//Main function for the WASM browser-side app to start up.
func main() {
	app.Route("/", new(Session))
	app.RunWhenOnBrowser()
}
