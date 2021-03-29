package main

import (

)

var CACHE Cache //in-memory cache of settings between all pages (cache.go)
var SC SessionCache //internal cache of session variables (formerly in the Session struct itself)

//Session Cache for the WASM binary
type Cache struct {
	LoggedIn bool
	PageTitle string

}

type SessionCache struct {
	//Main Page Information
	Current_header DialogPage
	Current_page MainPage
	Current_footer DialogPage

	//Pull-out panel information
	Panel_show  bool
	Panel_page DialogPage
	Panel_title string
	Panel_icon string

	//Popup information
	ShowPopup bool
	PopupText string
	PopupIcon string
	PopupYesNo PopupResult
	PopupString PopupStringResult
	Popup_page DialogPage
}
