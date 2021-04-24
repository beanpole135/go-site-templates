package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"strconv"
	"math/rand"
)

// Example/Template file for a "Dialog" component that can be used either in a popup
// or in a pull-out panel

type DialogExample struct {
	app.Compo
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
	color string
}

func (D *DialogExample) Render() app.HTMLDiv {
	if D.color == "" { D.color = "#aaaaaa" }
	return app.Div().Body(
		app.P().Text("I am an example component!\nI can be placed within all sorts of locations in the UI, including in pull-out panels or popup boxes!"),
		app.Button().Text("Random Color").OnClick(D.RandomColor).Style("background", D.color),
		app.Button().Text("Popup Dialog").OnClick(D.ShowPopup),
	)
}

func (D *DialogExample) ShowPopup(ctx app.Context, ev app.Event) {
	SESSION.Popup("icon", "You can even open popups from within components (this will close any previous popup though)!")
}

func (D *DialogExample) RandomColor(ctx app.Context, ev app.Event) {
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)
	D.color = "rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")"
}
