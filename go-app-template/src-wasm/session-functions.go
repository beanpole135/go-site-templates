package main

import (
	"time"
	"strings"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Various interaction functions for the global session

func (P *Session) ChangePage(pageid string, title string){
  if !CACHE.LoggedIn {
    SC.Current_page = &PageLogin{} //Use a login page here of some kind
  }else{
	CACHE.PageTitle = title
	switch pageid {
		default:
			SC.Current_page = &PageExample{}
	}
  }
  SC.CurPage = strings.TrimSuffix( strings.TrimPrefix(pageid,"/"), "/")
  //Ensure the slide-out panel and popups are cleared on page changes
  P.HidePopup()
  P.HidePanel()
  //P.Update()
}

// === Popup Dialogs ===
func (P *Session) HidePopup() {
	if SC.ShowPopup == false { return }
	SC.ShowPopup = false
	SC.Popup_page = nil
	SC.PopupYesNo = nil
	SC.PopupString = nil
	SC.Popup_menu = nil
	SC.Popup_pix_X = 0
	SC.Popup_pix_Y = 0
	//P.Update()
}

func (P *Session) Popup( icon string, text string){
	//Quick 3-second popup message
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupYesNo = nil
	SC.PopupString = nil
	//P.Update()
	go func(){
		time.Sleep(3 * time.Second)
		P.HidePopup()
	}()
}

func (P *Session) PopupTextBox( icon string, text string){
	// Popup message without a timed close (user has to click a button to make it go away)
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupYesNo = nil
	SC.PopupString = nil
	//P.Update()
}

func (P *Session) PopupYesNoBox( icon string, text string, callback PopupResult){
	// Popup with a Yes/No question
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupYesNo = callback
	//P.Update()
}

func (P *Session) PopupStringQuestion( icon string, text string, callback PopupStringResult) {
	// Popup requesting the user to type in some text
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupString = callback
	//P.Update()
}

func (P *Session) PopupDialog(icon string, text string, body DialogPage){
	// Generic popup with a custom render item for the body.
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.Popup_page = body
	//P.Update()
}

// === Context Menu functions ===
func (P *Session) PopupContextMenu(list []MenuItem, callback PopupStringResult, ctx app.Context){
	//Note: Use the "HidePopup()" function to cancel a context menu
	if ctx == nil {
		//Use the current mouse position
		SC.Popup_pix_X, SC.Popup_pix_Y = app.Window().CursorPosition()
	} else {
		//Use the bottom-left of the element clicked
		rect := ctx.JSSrc().Call("getBoundingClientRect")
		yoffset := app.Window().Get("pageYOffset").Int()
		xoffset := app.Window().Get("pageXOffset").Int()
		//put in absolute coords
		SC.Popup_pix_Y = rect.Get("bottom").Int() + yoffset
		SC.Popup_pix_X = rect.Get("left").Int() + xoffset
	}
	SC.Popup_menu = list
	SC.PopupString = callback
	SC.ShowPopup = true
	//P.Update()
}

// === Slide-out Panel functions ===
func (P *Session) ShowPanel( icon string, title string, body DialogPage) {
	// Slide-out panel to show a generic render item
	SC.Panel_page = body
	SC.Panel_title = title
	SC.Panel_icon = icon
	SC.Panel_show = true
	//P.Update()
}

func (P *Session) HidePanel(){
	if SC.Panel_show == false { return } //nothing to do
	SC.Panel_show = false
	SC.Panel_page = nil
	//P.Update()
}
