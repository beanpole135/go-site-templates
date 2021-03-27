package main

import (
	"time"
	"fmt"
)

// Various interaction functions for the global session

func (P *Session) ChangePage(pageid string, title string){
  fmt.Println("Change Page:", pageid, title, CACHE.LoggedIn)
  if !CACHE.LoggedIn {
    P.Current_page = &PageLogin{} //Use a login page here of some kind

  }else{
	CACHE.PageTitle = title
	switch pageid {
		default:
			P.Current_page = &PageExample{}
	}
  }
  fmt.Println(" - Page Changed: ", P.Current_page != nil)
  //Ensure the slide-out panel and popups are cleared on page changes
  P.HidePopup()
  P.HidePanel()
}

// === Popup Dialogs ===
func (P *Session) HidePopup() {
	if P.ShowPopup == false { return }
	P.ShowPopup = false
	P.Popup_page = nil
	P.PopupYesNo = nil
	P.PopupString = nil
	P.Update()
}

func (P *Session) Popup( icon string, text string){
	//Quick 3-second popup message
	P.PopupIcon = icon
	P.PopupText = text
	P.ShowPopup = true
	P.PopupYesNo = nil
	P.PopupString = nil
	P.Update()
	go func(){
		time.Sleep(3 * time.Second)
		P.HidePopup()
	}()
}

func (P *Session) PopupTextBox( icon string, text string){
	// Popup message without a timed close (user has to click a button to make it go away)
	P.PopupIcon = icon
	P.PopupText = text
	P.ShowPopup = true
	P.PopupYesNo = nil
	P.PopupString = nil
	P.Update()
}

func (P *Session) PopupYesNoBox( icon string, text string, callback PopupResult){
	// Popup with a Yes/No question
	P.PopupIcon = icon
	P.PopupText = text
	P.ShowPopup = true
	P.PopupYesNo = callback
	P.Update()
}

func (P *Session) PopupStringQuestion( icon string, text string, callback PopupStringResult) {
	// Popup requesting the user to type in some text
	P.PopupIcon = icon
	P.PopupText = text
	P.ShowPopup = true
	P.PopupString = callback
	P.Update()
}

func (P *Session) PopupDialog(icon string, text string, body DialogPage){
	// Generic popup with a custom render item for the body.
	P.PopupIcon = icon
	P.PopupText = text
	P.ShowPopup = true
	P.Popup_page = body
	P.Update()
}

// === Slide-out Panel functions ===
func (P *Session) ShowPanel( icon string, title string, body DialogPage) {
	// Slide-out panel to show a generic render item
	P.Panel_page = body
	P.Panel_title = title
	P.Panel_icon = icon
	P.Panel_show = true
	P.Update()
}

func (P *Session) HidePanel(){
	if P.Panel_show == false { return } //nothing to do
	P.Panel_show = false
	P.Panel_page = nil
	P.Update()
}
