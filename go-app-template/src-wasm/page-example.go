package main

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Example/Template file for a new page

type PageExample struct {
	app.Compo
	// Information can be stored here that only persists for the duration of this page
	// For persistent info storage, use the global CACHE object
	title *string //Link this over to the CACHE to get a global value
}

func (P *PageExample) Render() app.HTMLMain { //return is important for the "MainPage" interface definition
	fmt.Println("Render Page Example")
	if P.title == nil {
		P.title = &CACHE.PageTitle
	}
	return app.Main().Body(
		app.H1().Text(*P.title),
		app.H1().Text("Header1"),
		app.H2().Text("Header2"),
		app.H3().Text("Header3"),
		app.H4().Text("Header4"),
		app.H5().Text("Header5"),
		app.P().Text("I am a very pretty paragraph."),
		app.Span().Body(
			app.Button().Text("Test Popup").OnClick(P.ShowPopup),
			app.Button().Text("Test Menu").OnClick(P.ShowContextMenu),
			app.Button().Text("Test Panel").OnClick(P.ShowPanel),
		).Style("display", "flex").Style("justify-content", "space-between"),
	).Style("padding", "2em")
}

func (P *PageExample) ShowPopup(ctx app.Context, ev app.Event) {
	SESSION.Popup("icon", "Text string")
}

func (P *PageExample) ShowContextMenu(ctx app.Context, ev app.Event) {
	var menu []MenuItem
	menu = append(menu, MenuItem{ID: "1", Text: "Item 1"})
	menu = append(menu, MenuItem{ID: "2", Text: "Item 2"})
	menu = append(menu, MenuItem{ID: "3", Text: "Item 3"})
	menu = append(menu, MenuItem{ID: "4", Text: "Item 4"})
	//Align the context menu to the bottom of the element clicked
	SESSION.PopupContextMenu(menu, P.ContextMenuCallback, ctx)
	//Align the context menu to the current mouse position
	//SESSION.PopupContextMenu(menu, P.ContextMenuCallback, nil)
}

func (P *PageExample) ContextMenuCallback(id string) {
	fmt.Println("Got Context Menu Callback:", id)
	SESSION.Popup("", "Context Menu Item Selected:\n\n"+id)
}

func (P *PageExample) ShowPanel(ctx app.Context, ev app.Event) {
	dlg := new(DialogExample)
	SESSION.ShowPanel("icon", "Example", dlg)
}
