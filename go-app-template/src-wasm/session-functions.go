package main

import (
	"time"
)

// Various interaction functions for the global session


// === Popup Dialogs ===
func (P *Session) HidePopup() {
	P.showPopup = false
	P.popup_page = nil
	P.popupYesNo = nil
	P.popupString = nil
	P.Update()
}

func (P *Session) Popup( icon string, text string){
	//Quick 3-second popup message
	P.popupIcon = icon
	P.popupText = text
	P.showPopup = true
	P.popupYesNo = nil
	P.popupString = nil
	P.Update()
	go func(){
		time.Sleep(3 * time.Second)
		P.HidePopup()
	}()
}

func (P *Session) PopupText( icon string, text string){
	// Popup message without a timed close (user has to click a button to make it go away)
	P.popupIcon = icon
	P.popupText = text
	P.showPopup = true
	P.popupYesNo = nil
	P.popupString = nil
	P.Update()
}

func (P *Session) PopupYesNo( icon string, text string, callback popupResult){
	// Popup with a Yes/No question
	P.popupIcon = icon
	P.popupText = text
	P.showPopup = true
	P.popupYesNo = callback
	P.Update()
}

func (P *Session) PopupStringQuestion( icon string, text string, callback popupStringResult) {
	// Popup requesting the user to type in some text
	P.popupIcon = icon
	P.popupText = text
	P.showPopup = true
	P.popupString = callback
	P.Update()
}

func (P *Session) PopupDialog(icon string, text string, body DialogPage){
	// Generic popup with a custom render item for the body.
	P.popupIcon = icon
	P.popupText = text
	P.showPopup = true
	P.popup_page = body
	P.Update()
}

// === Slide-out Panel functions ===
func (P *Session) ShowPanel( icon string, title string, body DialogPage) {
	// Slide-out panel to show a generic render item
	P.panel_page = body
	P.panel_title = title
	P.panel_icon = icon
	P.panel_show = true
	P.Update()
}

func (P *Session) HidePanel(){
	P.panel_show = false
	P.panel_page = nil
	P.Update()
}
