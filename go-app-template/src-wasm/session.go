package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"strconv"
)

// ====================
// The "Session" is the main top-level component for the app.
//  and provides site-wide functionality for all pages/components to use as needed
// - See session-function.go for the list of available popups/panels/etc.
// - Use "SESSION.Update" to trigger a refresh of the UI after something changes programmatically
// ====================
// Due to how go-app v8 initializes things, there is a separate SessionCache (SC)
// structure for handling state-fields of the session. This is not typically used directly
// but is also available globally should you need it.
// ====================
var SESSION *Session //Primary page-render system (session.go)
var SC SessionCache  //internal cache of session variables (formerly in the Session struct itself)

// Primary Page-type definitions
type MainPage interface {
	Render() app.HTMLMain
}

type DialogPage interface {
	Render() app.HTMLDiv
}

type MenuItem struct {
	ID   string
	Icon string
	Text string
}

// Interaction functions for popup/dialog callbacks
type PopupResult func(bool)
type PopupStringResult func(string)

// Internal Session Cache (global session objects - set dynamically in session-functions.go)
type SessionCache struct {
	//Main Page Information
	CurPage        string
	Current_header DialogPage
	Current_page   MainPage
	Current_footer DialogPage

	//Pull-out panel information
	Panel_show  bool
	Panel_page  DialogPage
	Panel_title string
	Panel_icon  string

	//Popup information
	ShowPopup   bool
	PopupText   string
	PopupIcon   string
	PopupYesNo  PopupResult
	PopupString PopupStringResult
	Popup_page  DialogPage
	Popup_menu  []MenuItem
	Popup_pix_X int
	Popup_pix_Y int
}

// Primary Visual Element for the UI
// This sets up all the individual global-visuals and shows the right pages
type Session struct {
	app.Compo
}

// ==== MAIN RENDER ROUTINE ====
func (S *Session) Render() app.UI {
	if SESSION != S {
		SESSION = S
		SESSION.ChangePage("/", "Dashboard")
	}
	if SC.Current_page == nil {
		return app.Div()
	}
	return app.Div().Body(
		S.RenderHeader(),
		SC.Current_page.Render(),
		S.RenderFooter(),
		S.RenderPanel(),
		S.RenderPopup(),
		S.RenderCoverWindow(),
	)
}

// ==== Render routines for individual global components ====
func (S *Session) RenderHeader() app.UI {
	if SC.Current_header == nil {
		return nil
	}
	return SC.Current_header.Render().Class("body-header")
}
func (S *Session) RenderFooter() app.UI {
	if SC.Current_footer == nil {
		return nil
	}
	return SC.Current_footer.Render().Class("body-footer")
}

func (P *Session) RenderPopup() app.UI {
	if !SC.ShowPopup {
		return nil
	}
	dlg := P.RenderPopupDialog()
	if SC.Popup_page != nil {
		dlg.Body(
			app.Span().Class("align-row").Body(
				app.P().Text(SC.PopupText),
				app.Button().Text("cancel").OnClick(P.HidePopupCallback),
			),
			SC.Popup_page.Render(),
		).Style("display","flex").Style("flex-direction","column")
	} else if SC.PopupYesNo != nil {
		dlg.Body(
			app.P().Text(SC.PopupText).Style("font-size", "large"),
			app.Div().Body(
				app.Button().ID("no").Text("No").OnClick(P.PopupAnswer),
				app.Button().ID("yes").Text("Yes").OnClick(P.PopupAnswer),
			).
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "space-evenly"),
		).Style("display","flex").Style("flex-direction","column")
	} else if SC.Popup_menu != nil {
		dlg.Body(
			app.Range(SC.Popup_menu).Slice(P.RenderMenuItem),
		).Style("display","flex").Style("flex-direction","column")

	} else if SC.PopupString != nil {
		dlg.Body(
			app.P().Text(SC.PopupText).Style("font-size", "large"),
			app.Input().Type("text").ID("dialog_text_input").Style("margin-bottom", "1ex"),
			app.Div().Body(
				app.Button().ID("no").Text("Cancel").OnClick(P.PopupAnswer),
				app.Button().ID("yes").Text("Continue").OnClick(P.PopupAnswer),
			).Style("display", "flex").Style("align-items", "center").Style("justify-content", "space-evenly"),
		).Style("display","flex").Style("flex-direction","column")
	} else {
		dlg.Body(
			app.P().Text(SC.PopupText),
			app.Button().Text("Continue").OnClick(P.PopupAnswer),
		)
	}
	if SC.Popup_pix_X != 0 && SC.Popup_pix_Y != 0 {
		maxX, maxY := app.Window().Size()
		maxX = maxX - SC.Popup_pix_X
		maxY = maxY - SC.Popup_pix_Y
		dlg.
		Style("top", strconv.Itoa(SC.Popup_pix_Y)+"px" ).
		Style("left", strconv.Itoa(SC.Popup_pix_X)+"px" ).
		Style("max-height", strconv.Itoa(maxY)+"px" ).
		Style("max-width", strconv.Itoa(maxX)+"px" ).
		Style("margin","0")
	}
	return dlg
}
func (P *Session) RenderPopupDialog() app.HTMLDialog {
	return app.Dialog().Hidden(false).Open(true).
		Class("popup").
		Style("border", "0.2ex solid var(--COLOR-Accent)").
		Style("background", "var(--COLOR-Background-dark)").
		Style("color", "var(--COLOR-text-light)").
		Style("position", "absolute").
		Style("top", "40%").
		Style("border-radius", "1ex").
		Style("z-index", "1000").
		Style("max-width", "50%")
}
func (P *Session) RenderCoverWindow() app.UI {
	if !SC.ShowPopup && SC.Panel_page == nil{
		return nil
	}
	//Need to show a cover window
	cover := app.Div().
		Style("background","#00000010").
		Style("position","absolute").
		Style("top","0").
		Style("bottom","0").
		Style("left","0").
		Style("right","0")
	if SC.ShowPopup{
		return cover.Style("z-index","999").OnClick(P.HidePopupCallback)
	}else{
		return cover.Style("z-index","499").OnClick(P.HidePanelCallback)
	}
		
}

func (P *Session) RenderMenuItem(index int) app.UI {
	M := SC.Popup_menu[index]
	if M.Icon == "" {
		return app.Span().Class("menuitem").ID(M.ID).Body(
			app.Text(M.Text),
		).OnClick(P.PopupAnswer)

	} else {
		return app.Span().Class("menuitem").ID(M.ID).Body(
			app.Text(M.Text),
		).OnClick(P.PopupAnswer)
	}
}

/// ===  Panel Pullout ====

func (P *Session) RenderPanel() app.UI {
	if SC.Panel_page != nil {
		return app.Dialog().Class("panel").ID("right-panel").Hidden(!SC.Panel_show).Open(SC.Panel_show).Body(
				app.Span().Class("panel-header").Body(
					app.H2().Text(SC.Panel_title),
					app.Button().Text("Cancel").OnClick(P.HidePanelCallback).Style("padding", "1ex"),
				),
				SC.Panel_page.Render().Class("panel-content"),
			).
			Style("width", "50%").
			Style("z-index", "500")
	} else {
		return app.Dialog().Class("panel").ID("right-panel").Hidden(!SC.Panel_show).Open(SC.Panel_show).Body(
			app.Span().Class("panel-header").Body(
				app.H2().Text(SC.Panel_title),
				app.Button().Text("Cancel").OnClick(P.HidePanelCallback).Style("padding", "1ex"),
			),
			app.Div().Class("panel-content"),
		).
		Style("width", "0").
		Style("z-index", "500")
	}
}

// ==== CALLBACK FUNCTIONS ====
// These are used by the global components internally - not generally used elsewhere
func (P *Session) HidePanelCallback(ctx app.Context, e app.Event) {
	P.HidePanel()
}

func (P *Session) HidePopupCallback(ctx app.Context, e app.Event) {
	P.HidePopup()
}

func (P *Session) PopupAnswer(ctx app.Context, e app.Event) {
	//Note: We need to "HidePopup" right before the callback, just in case the callback
	// Triggers another popup of some kind.
	if SC.PopupYesNo != nil {
		callback := SC.PopupYesNo
		id := ctx.JSSrc.Get("id").String()
		P.HidePopup()
		switch id {
		case "no":
			callback(false)
		case "yes":
			callback(true)
		}
	} else if SC.PopupString != nil {
		callback := SC.PopupString
		id := ctx.JSSrc.Get("id").String()
		if id == "no" {
			//cancelled - do nothing
			P.HidePopup()
		} else {
			//Need to read the string value from the input box
			textinput := app.Window().GetElementByID("dialog_text_input")
			if textinput.Truthy() {
				//Special text-input box was used in this popup
				txt := textinput.Get("value").String()
				P.HidePopup()
				callback(txt)
			} else {
				//Generic item selected - send the ID back
				P.HidePopup()
				callback(id)
			}
		}
	} else {
		P.HidePopup()
	}

}
