package main

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Example/Template file for a new page

type PageLogin struct {
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
}

func (P *PageLogin) Render() app.HTMLMain { //return is important for the "MainPage" interface definition
	fmt.Println("Render PageLogin")
	return app.Main().Body(
		app.Button().Text("Login Test").OnClick(P.Login),
	)
}

func (P *PageLogin) Login(ctx app.Context, ev app.Event) {
	err := SendAPI_Login("test", "test")
	if err == nil {
		CACHE.LoggedIn = true
		SESSION.ChangePage("/", "")
		SESSION.Popup("icon", "Successfully logged in!")
	} else {
		SESSION.Popup("icon", "Error Logging in...")
	}
}
