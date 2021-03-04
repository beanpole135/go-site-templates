package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)
// GLOBALS for all pages
var CACHE Cache //in-memory cache of settings between all pages (cache.go)
var SESSION Session //Primary page-render system (session.go)

func main() {
	app.RouteWithRegexp("/", &SESSION)
	app.Run()
}
