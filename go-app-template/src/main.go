package main

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"net/http"
	"os"
)

// ===================
// GLOBAL VARIABLES
// ===================
var progname string = "ChangeME"
var SETTINGS Settings

// ===================

func main() {
	//SEO Rendering functionality
	app.Route("/", new(SEO))
	app.RunWhenOnBrowser() //ignored on the server-side builds
	//For the server-side, everything below gets used
	if len(os.Args) > 1 {
		SETTINGS = readSettings(os.Args[1])
	} else {
		SETTINGS = readSettings("")
	}
	SetupLoginAuth()
	//Setup the go-app system
	h := &app.Handler{
		Title:           progname,
		Author:          progname,
		Name:            progname,
		ShortName:       progname,
		Icon:            app.Icon{Default: "/web/favicon.png"},
		ThemeColor:      "#111111",
		BackgroundColor: "#111111",
		Styles:          []string{"/web/style.css"},
		LoadingLabel:    "Please wait",
	}
	//Static File/Dir handling
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(InstallDir+"/web/static"))))
	http.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(http.Dir(InstallDir+"/web"))))
	//Interactive API handling
	http.HandleFunc("/api/", HandleAPI)                 //api.go
	http.HandleFunc("/login-oauth/", StartOAuthLogin)   //auth_oauth.go
	http.HandleFunc("/handle-oauth/", HandleOAuthLogin) //auth_oauth.go
	//The routing specific to loading the WASM app
	http.HandleFunc("/", h.ServeHTTP)
	//Now start listening
	fmt.Println(" - listening:", SETTINGS.ListenURL)
	http.ListenAndServe(SETTINGS.ListenURL, nil)
}
