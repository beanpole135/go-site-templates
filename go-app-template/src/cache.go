package main

import (

)

var CACHE Cache //in-memory cache of settings between all pages (cache.go)

//Session Cache for the WASM binary
type Cache struct {
	LoggedIn bool
	PageTitle string
}
