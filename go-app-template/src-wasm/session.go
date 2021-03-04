package main

import (
	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

// Primary Visual Element for the UI
// This sets up all the individual global-visuals and shows the right pages
type Session struct {
	app.Compo
	//Main Page Information
	current_header DialogPage
	current_page MainPage
	current_footer DialogPage

	//Pull-out panel information
	panel_show  bool
	panel_page DialogPage
	panel_title string
	panel_icon string

	//Popup information
	showPopup bool
	popupText string
	popupIcon string
	popupYesNo popupResult
	popupString popupStringResult
	popup_page DialogPage

}

// Primary Page-type definitions
type MainPage interface {
	Render() app.HTMLMain
}

type DialogPage interface {
	Render() app.HTMLDiv
}

// Interaction functions for popup/dialog callbacks
type popupResult func(bool)
type popupStringResult func(string)

// ==== MAIN RENDER ROUTINE ====
func (S *Session) Render() app.UI {
	return app.Main().Body(
		S.RenderHeader(),
		app.P().Text("I am alive!!!"),
		S.RenderFooter(),
		S.RenderPanel(),
		S.RenderPopup(),
	)
}

// ==== Render routines for individual global components ====
func (S *Session) RenderHeader() app.UI {
	if S.current_header == nil { return nil }
	return S.current_header.Render().Class("body-header")
}
func (S *Session) RenderFooter() app.UI {
	if S.current_footer == nil { return nil }
	return S.current_footer.Render().Class("body-footer")
}

func (P *Session) RenderPopup() app.UI {
	if P.showPopup {
		if P.popup_page != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
					app.Span().Class("align-row").Body(
						app.P().Text( P.popupText ),
						app.Button().Text("cancel").OnClick(P.HidePopupCallback),
					),
					P.popup_page.Render(),
				).
				Style("border","1ex solid var(--COLOR-Accent)").
				Style("background","var(--COLOR-Background-dark)").
				Style("color","var(--COLOR-text-light)").
				Style("position","absolute").
				Style("top","40%").
				Style("border-radius","1ex").
				Style("z-index","1000").
				Style("max-width","75%")
		} else if P.popupYesNo != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
				app.P().Text(P.popupText).Style("font-size","large"),
				app.Div().Body(
					app.Button().ID("no").Text("No").OnClick(P.popupAnswer),
					app.Button().ID("yes").Text("Yes").OnClick(P.popupAnswer),
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
		} else if P.popupString != nil {
			return app.Dialog().Hidden(false).Open(true).Body( 
				app.P().Text(P.popupText).Style("font-size","large"),
				app.Input().Type("text").ID("dialog_text_input").Style("margin-bottom","1ex"),
				app.Div().Body(
					app.Button().ID("no").Text("Cancel").OnClick(P.popupAnswer),
					app.Button().ID("yes").Text("Continue").OnClick(P.popupAnswer),
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
					app.P().Text(P.popupText),
					app.Button().Text("Continue").OnClick(P.popupAnswer),
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
  if P.panel_page != nil {
	return app.Dialog().Class("panel").ID("right-panel").Hidden(!P.panel_show).Open(P.panel_show).Body(
		app.Span().Class("panel-header").Body( 
			app.H2().Text(P.panel_title),
			app.Button().Text("Cancel").OnClick(P.HidePanelCallback).Style("padding","1ex"),
		),
		P.panel_page.Render().Class("panel-content"),
	).
	Style("width","50%").
	Style("z-index","500")
  }else{
	return app.Dialog().Class("panel").ID("right-panel").Hidden(!P.panel_show).Open(P.panel_show).Body(
		app.Span().Class("panel-header").Body( 
			app.H2().Text(P.panel_title),
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

func (P *Session) popupAnswer(ctx app.Context, e app.Event){
	P.showPopup = false
	P.popupText = ""
	P.popupIcon = ""
	if P.popupYesNo != nil {
		id := ctx.JSSrc.Get("id").String()
		switch id {
			case "no": P.popupYesNo(false)
			case "yes": P.popupYesNo(true)	
		}
		P.popupYesNo = nil
	}else if P.popupString != nil {
		id := ctx.JSSrc.Get("id").String()
		if id == "no" {
			//cancelled - do nothing
			//P.popupString("")
		} else {
			//Need to read the string value from the input box
			input := app.Window().GetElementByID("dialog_text_input").Get("value").String()
			P.popupString(input)
		}
		P.popupString = nil
	}
	P.HidePopup()
}
