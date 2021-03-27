package main

import (
	"fmt"
	"encoding/json"
	"strings"
	"errors"
	"net/http"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"bytes"
	"io/ioutil"
)

type api_request struct {
	API string
	ArgData []byte
	ReplyData interface{}
	UserIDAuthorized string
	SessionToken string
}

func (R *api_request) Evaluate() error {
	//General process
	// 1. Switch over the API string
	// 2. Convert the ArgData into a structure for the designated API (typically using JSON)
	// 3. Run a function to "do the thing" and get a generic reply/error
	// 4. Put the reply into the "ReplyData" field, and return the error 
	err := errors.New("Unknown API")
	switch R.API {
		case "do-something":
			//R.ReplyData, err = DoAPISomething(UserIDAuthorized, R.ArgData)
	}
	return err
}

func HandleAPI(w http.ResponseWriter, r *http.Request){
	//This is setup for an HTTP/REST request/reply format. Not a websocket format
	//Parse Input body / URL
	var A api_request
	A.API = strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/api/"), "/") // "host.com/api/test/something/" -> API: "test/something"
	
	A.UserIDAuthorized, A.SessionToken = CheckLogin(w,r)
	// UserIDAuthorized is non-empty if valid credentials were provided (either a token or basic auth user/pass)
	// SessionToken is non-empty if they have an active "session" token right now.
	err := errors.New("Unauthorized")
	//Check for authorized access or specific login/logout API's
	if A.UserIDAuthorized == "" {
		w.WriteHeader(401) //unauthorized
		return
	} else {
		//Evaluate Request
		if A.API =="login" {
			//Generate a login token and save it as a cookie in their browser
			if A.SessionToken == "" {
				GenerateLoginToken(w, r, A.UserIDAuthorized)
			}
			err = nil
		}else if A.API == "logout" {
			if A.SessionToken != "" {
				RemoveLoginToken(w, r, A.SessionToken)
			}
			err = nil
		}else{
			err = A.Evaluate()
		}
	}

	//Send Reply
	if err != nil {
		//Bad Request
		fmt.Println("Bad Request:", err, "\n"+A.API+" :", string(A.ArgData) )
		w.WriteHeader(400)
	}else{
		//Good Request
		if A.ReplyData == nil { 
			//Nothing to send back - just return the 200-OK header
			w.WriteHeader(200)
		}else{
			//Send back the reply as JSON data
			dat, _ := json.Marshal(A.ReplyData)
			w.Header().Set("Content-Type", "application/json")
			w.Write(dat)
		}
	}
}

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
