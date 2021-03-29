package main

import (
	"time"
	"fmt"
)

// Various interaction functions for the global session

func (P *Session) ChangePage(pageid string, title string){
  fmt.Println("Change Page:", pageid, title, CACHE.LoggedIn)
  if !CACHE.LoggedIn {
    SC.Current_page = &PageLogin{} //Use a login page here of some kind

  }else{
	CACHE.PageTitle = title
	switch pageid {
		default:
			SC.Current_page = &PageExample{}
	}
  }
  fmt.Println(" - Page Changed: ", SC.Current_page != nil)
  //Ensure the slide-out panel and popups are cleared on page changes
  P.HidePopup()
  P.HidePanel()
}

// === Popup Dialogs ===
func (P *Session) HidePopup() {
	if SC.ShowPopup == false { return }
	SC.ShowPopup = false
	SC.Popup_page = nil
	SC.PopupYesNo = nil
	SC.PopupString = nil
	P.Update()
}

func (P *Session) Popup( icon string, text string){
	//Quick 3-second popup message
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupYesNo = nil
	SC.PopupString = nil
	P.Update()
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
	P.Update()
}

func (P *Session) PopupYesNoBox( icon string, text string, callback PopupResult){
	// Popup with a Yes/No question
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupYesNo = callback
	P.Update()
}

func (P *Session) PopupStringQuestion( icon string, text string, callback PopupStringResult) {
	// Popup requesting the user to type in some text
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.PopupString = callback
	P.Update()
}

func (P *Session) PopupDialog(icon string, text string, body DialogPage){
	// Generic popup with a custom render item for the body.
	SC.PopupIcon = icon
	SC.PopupText = text
	SC.ShowPopup = true
	SC.Popup_page = body
	P.Update()
}

// === Slide-out Panel functions ===
func (P *Session) ShowPanel( icon string, title string, body DialogPage) {
	// Slide-out panel to show a generic render item
	SC.Panel_page = body
	SC.Panel_title = title
	SC.Panel_icon = icon
	SC.Panel_show = true
	P.Update()
}

func (P *Session) HidePanel(){
	if SC.Panel_show == false { return } //nothing to do
	SC.Panel_show = false
	SC.Panel_page = nil
	P.Update()
}
