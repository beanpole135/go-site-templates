package main

import (
	"time"
)

var CACHE Cache //in-memory cache of settings between all pages (cache.go)

//Session Cache for the WASM binary
type Cache struct {
	LoggedIn  bool
	PageTitle string
}

func (C *Cache) RandomizeOnTimer() {
	//This is a small function to simulate the app having the data cache updated via an external source
	// Such as a websocket event, asynchronous fetch/timer, etc
	// Changes made here should instantly reflect within the current app page
	go func() {
		for now := range time.Tick(time.Second) {
			C.PageTitle = now.Format("Mon Jan 2, 2006, 15:04:05 MST")
			if (now.Second() % 2) == 0 {
				//Only force-update the session every other second.
				// If the go-app system is detecting the string change, then we should see it keep time every second.
				// otherwise, we will only see it change every other second.
				SESSION.Update()
			}
		}
	}()
}
