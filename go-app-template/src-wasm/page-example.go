package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"fmt"
)

// Example/Template file for a new page

type PageExample struct {
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
}

func (P *PageExample) Render() app.HTMLMain { //return is important for the "MainPage" interface definition
	fmt.Println("Render Page Example")
	return app.Main().Body(
		app.H1().Text("Header1"),
		app.H2().Text("Header2"),
		app.H3().Text("Header3"),
		app.H4().Text("Header4"),
		app.H5().Text("Header5"),
		app.P().Text("I am a very pretty paragraph."),
		app.Button().Text("Test Popup").OnClick(P.ShowPopup),
	)
}

func (P *PageExample) ShowPopup(ctx app.Context, ev app.Event){
	SESSION.Popup("icon","Text string")
}
