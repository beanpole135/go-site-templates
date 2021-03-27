package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"fmt"
)

//var SESSION Session //Primary page-render system (session.go)

// Primary Visual Element for the UI
// This sets up all the individual global-visuals and shows the right pages
type Session struct {
	app.Compo
}

// Primary Page-type definitions
type MainPage interface {
	Render() app.HTMLMain
}

type DialogPage interface {
	Render() app.HTMLDiv
}

// Interaction functions for popup/dialog callbacks
type PopupResult func(bool)
type PopupStringResult func(string)

// ==== MAIN RENDER ROUTINE ====
func (S *Session) Render() app.UI {
	fmt.Println("Session Render Page:", SC.Current_page != nil )
	if SC.Current_page == nil {
		return app.Div()
	}
	return app.Div().Body(
		S.RenderHeader(),
		SC.Current_page.Render(),
		S.RenderFooter(),
		S.RenderPanel(),
		S.RenderPopup(),
	)
}

// ==== Render routines for individual global components ====
func (S *Session) RenderHeader() app.UI {
	if SC.Current_header == nil { return nil }
	return SC.Current_header.Render().Class("body-header")
}
func (S *Session) RenderFooter() app.UI {
	if SC.Current_footer == nil { return nil }
	return SC.Current_footer.Render().Class("body-footer")
}

func (P *Session) RenderPopup() app.UI {
	if SC.ShowPopup {
		if SC.Popup_page != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
					app.Span().Class("align-row").Body(
						app.P().Text( SC.PopupText ),
						app.Button().Text("cancel").OnClick(P.HidePopupCallback),
					),
					SC.Popup_page.Render(),
				).
				Style("border","1ex solid var(--COLOR-Accent)").
				Style("background","var(--COLOR-Background-dark)").
				Style("color","var(--COLOR-text-light)").
				Style("position","absolute").
				Style("top","40%").
				Style("border-radius","1ex").
				Style("z-index","1000").
				Style("max-width","75%")
		} else if SC.PopupYesNo != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
				app.P().Text(SC.PopupText).Style("font-size","large"),
				app.Div().Body(
					app.Button().ID("no").Text("No").OnClick(P.PopupAnswer),
					app.Button().ID("yes").Text("Yes").OnClick(P.PopupAnswer),
				).
				Style("display","flex").
				Style("align-items","center").
				Style("justify-content","space-evenly"),
			).
			Style("border","1ex solid var(--COLOR-Accent)").
			Style("background","var(--COLOR-Background-dark)").
			Style("color","var(--COLOR-text-light)").
			Style("position","absolute").
			Style("top","40%").
			Style("border-radius","1ex").
			Style("z-index","1000").
			Style("max-width","75%")
		} else if SC.PopupString != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
				app.P().Text(SC.PopupText).Style("font-size","large"),
				app.Input().Type("text").ID("dialog_text_input").Style("margin-bottom","1ex"),
				app.Div().Body(
					app.Button().ID("no").Text("Cancel").OnClick(P.PopupAnswer),
					app.Button().ID("yes").Text("Continue").OnClick(P.PopupAnswer),
				).Style("display","flex").Style("align-items","center").Style("justify-content","space-evenly"),
			).
			Style("border","1ex solid var(--COLOR-Accent)").
			Style("background","var(--COLOR-Background-dark)").
			Style("color","var(--COLOR-text-light)").
			Style("position","absolute").
			Style("top","40%").
			Style("border-radius","1ex").
			Style("z-index","1000").
			Style("display","flex").
			Style("flex-direction", "column").
			Style("max-width","75%")
		} else {
			return app.Dialog().Hidden(false).Open(true).Body( 
					app.P().Text(SC.PopupText),
					app.Button().Text("Continue").OnClick(P.PopupAnswer),
				).
				Style("border","1ex solid var(--COLOR-Accent)").
				Style("background","var(--COLOR-Background-dark)").
				Style("color","var(--COLOR-text-light)").
				Style("position","absolute").
				Style("top","40%").
				Style("border-radius","1ex").
				Style("z-index","1000").
				Style("max-width","75%")
		}
	}
	return nil
}

/// ===  Panel Pullout ====

func (P *Session) RenderPanel() app.UI {
  if SC.Panel_page != nil {
	return app.Dialog().Class("panel").ID("right-panel").Hidden(!SC.Panel_show).Open(SC.Panel_show).Body(
		app.Span().Class("panel-header").Body( 
			app.H2().Text(SC.Panel_title),
			app.Button().Text("Cancel").OnClick(P.HidePanelCallback).Style("padding","1ex"),
		),
		SC.Panel_page.Render().Class("panel-content"),
	).
	Style("width","50%").
	Style("z-index","500")
  }else{
	return app.Dialog().Class("panel").ID("right-panel").Hidden(!SC.Panel_show).Open(SC.Panel_show).Body(
		app.Span().Class("panel-header").Body( 
			app.H2().Text(SC.Panel_title),
			app.Button().Text("Cancel").OnClick(P.HidePanelCallback).Style("padding","1ex"),
		),
		app.Div().Class("panel-content"),
	).
	Style("width","0").
	Style("z-index","500")
  }
}

// ==== CALLBACK FUNCTIONS ====
// These are used by the global components internally - not generally used elsewhere
func (P *Session) HidePanelCallback(ctx app.Context, e app.Event){
	P.HidePanel()
}

func (P *Session) HidePopupCallback(ctx app.Context, e app.Event) {
	P.HidePopup()
}

func (P *Session) PopupAnswer(ctx app.Context, e app.Event){
	SC.ShowPopup = false
	SC.PopupText = ""
	SC.PopupIcon = ""
	if SC.PopupYesNo != nil {
		id := ctx.JSSrc.Get("id").String()
		switch id {
			case "no": SC.PopupYesNo(false)
			case "yes": SC.PopupYesNo(true)	
		}
		SC.PopupYesNo = nil
	}else if SC.PopupString != nil {
		id := ctx.JSSrc.Get("id").String()
		if id == "no" {
			//cancelled - do nothing
			//P.PopupString("")
		} else {
			//Need to read the string value from the input box
			input := app.Window().GetElementByID("dialog_text_input").Get("value").String()
			SC.PopupString(input)
		}
		SC.PopupString = nil
	}
	P.HidePopup()
}
