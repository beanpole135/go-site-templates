package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"math/rand"
	"strconv"
	"time"
)

// Example/Template file for a "Dialog" component that can be used either in a popup
// or in a pull-out panel

var tmp string = "#aaaaaa"

func (D *DialogExample) RandomColor2() {
	go func() {
		for now := range time.Tick(time.Second) {
			r := rand.Intn(255)
			g := rand.Intn(255)
			b := rand.Intn(255)
			color := "rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")"
			if now.Second()%2 == 0 {
				tmp = color
			} else {
				D.color = color
			}
		}
	}()
}

type DialogExample struct {
	app.Compo
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
	color      string
	color2     *string
	timestring *string
}

func (D *DialogExample) Render() app.HTMLDiv {
	if D.color == "" {
		D.color = "#aaaaaa"
	}
	if D.color2 == nil {
		D.color2 = &tmp
		D.RandomColor2()
	}
	if D.timestring == nil {
		D.timestring = &CACHE.PageTitle
	}
	return app.Div().Body(
		app.P().Text("I am an example component!\nI can be placed within all sorts of locations in the UI, including in pull-out panels or popup boxes!"),
		app.Button().Text("Random Color").OnClick(D.RandomColor).Style("background", D.color),
		app.Button().Text("Popup Dialog").OnClick(D.ShowPopup).Style("background", *D.color2),
		app.Text(*D.timestring),
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
