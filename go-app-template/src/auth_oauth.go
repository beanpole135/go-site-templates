package main

import (
	"net/http"
	"strings"
)

func StartOAuthLogin(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/login-oauth/"), "/") //matches main.go routing
	url := "/"
	switch provider {
	default:
		//TO-DO nothing specified here yet
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	//Note: When configuring the OAuth provider the callback URL needs to be: <this_host:port>/handle-oauth/<provider>
	provider := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/handle-oauth/"), "/") //matches main.go routing
	login_user := ""                                                                      //Set this to a user ID if the login was successful
	switch provider {
	default:
		//TO-DO nothing specified here yet
	}
	if login_user != "" {
		//Got a successful login - generate a login session/token for the user
		GenerateLoginToken(w, r, login_user)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther) //redirect back to the WASM site
}
