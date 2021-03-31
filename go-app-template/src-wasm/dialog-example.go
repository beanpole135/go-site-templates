package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

// Example/Template file for a "Dialog" component that can be used either in a popup
// or in a pull-out panel

type DialogExample struct {
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
}

func (D *DialogExample) Render() app.HTMLDiv {
	return app.Div().Body(
		app.P().Text("I am an example component!\nI can be placed within all sorts of locations in the UI, including in pull-out panels or popup boxes!"),
	)
}

func (D *DialogExample) ShowPopup(ctx app.Context, ev app.Event) {
	SESSION.Popup("icon", "You can even open popups from within components (this will close any previous popup though)!")
}
