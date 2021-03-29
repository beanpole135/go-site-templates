package main

import (
	"encoding/json"
	"net/http"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"bytes"
	"io/ioutil"
)

// Generic funtion for sending an arbitrary API call
// This should be inside functions for specific API calls
func SendAPI(api string, args interface{}) ([]byte, error) {
	//Send an API call to the backend and return the reply + ok state
	rawurl := app.Window().URL()
	rawurl.Path = "/api/"+api
	rawurl.RawQuery = ""
	rawurl.Fragment = ""
	app.Log("Send API:", rawurl.String() )
	var resp *http.Response
	var err error
	if args != nil {
		//Need to send some data as well
		dat, _ := json.Marshal(args)
		buf := bytes.NewBuffer(dat)
		req, _ := http.NewRequest("GET", rawurl.String(), buf)
		resp, err = http.DefaultClient.Do(req)
		//resp, err = http.Post(rawurl.String(), "application/json", buf) 
	}else {
		resp, err = http.Get(rawurl.String())
	}
	//Now read the reply and return
	if err != nil { return []byte{}, err } //error in request/reply
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func SendAPI_Login(username string, password string) error {
	//Special function to include the username/password as part of the login API
	// Does not wrap the SendAPI function like most will
	rawurl := app.Window().URL()
	rawurl.Path = "/login"
	rawurl.RawQuery = ""
	rawurl.Fragment = ""
	req, err := http.NewRequest("GET",rawurl.String(), nil)
	req.SetBasicAuth(username,password)
	_, err = http.DefaultClient.Do(req)
	return err
}

func SendAPI_Logout() error {
	//Standard form function - just uses the SendAPI function inside
	_, err := SendAPI("logout", nil)
	return err
}
