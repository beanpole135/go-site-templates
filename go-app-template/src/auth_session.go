package main

import (
	"net/http"
	"time"
	"github.com/gorilla/securecookie"
	"sync"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func PerformLogin(user string, password string) string {
	//Returns: userID if the login is successful

	// SETUP AS APPROPRIATE FOR APPLICATION
	// This could use LDAP, or an internal database/lookup, or a remote login call as needed
	fmt.Println("PerformLogin() function is not setup yet. Allowing all access")
	return user
}

// This file is the shim-layer between the OAuth login process and the client browser
// It manages a single cookie on the client browser, and uses that to allow temporary access
//   to loading the WASM Web App (after login) as well as the Websocket connection that the
//   web app establishes with the server after it loads in the browser

type AuthSession struct {
	Login_name string	`json:"login_name"`		//User ID / username
	Expires int64		`json:"-"`
}

func (A *AuthSession) SetExpire(temporary bool) {
	if temporary {
		//Quick temporary timeout (for loading previous sessions)
		A.Expires = time.Now().Add(time.Minute * 5).Unix()
	}else{
		//Long-term timeout for standard login
		A.Expires = time.Now().Add(time.Hour * 12).Unix()
	}
}

var secure *securecookie.SecureCookie = securecookie.New( securecookie.GenerateRandomKey(64), nil)
var loginSessions map[string]*AuthSession = make(map[string]*AuthSession)
var cookiename string = progname+"-session"
var lsesslock sync.Mutex
var sesschange bool = false

func SetupLoginAuth(){
	//Initialization function
	changed := false
	if SETTINGS.LSessHashKey == nil { SETTINGS.LSessHashKey = securecookie.GenerateRandomKey(64) ; changed = true }
	if SETTINGS.LSessBlockKey == nil { SETTINGS.LSessBlockKey = securecookie.GenerateRandomKey(32) ; changed = true }
	secure = securecookie.New( SETTINGS.LSessHashKey, nil)
	//Load last session hash
	if changed {
		SETTINGS.Save() //Save the new keys into the config for next time
	} else {
		//The cache only works if the hash/block keys are unchanged between runs
		readAuthSessionCache()
	}
	go PruneTokensThread()
}

func PruneTokensThread() {
	//Designed to be run in a separate go routine to do periodic cleanup of tokens
	for now := range time.Tick(time.Minute) {
		check := now.Unix()
		//Check all the auth sessions and wipe any that are expired
		for tok, auth := range loginSessions {
			if auth.Expires > 0 && auth.Expires < check {
				//Expired token - wipe it
				deleteAuthSession(tok)
			}
		}
		writeAuthSessionCache()
	}
}

// ============================
// Auth Session Token Management
// ============================
func saveAuthSession(cookieTok string, auth *AuthSession){
  if auth.Expires == 0 {
    auth.SetExpire(true)
  }
  lsesslock.Lock()
  loginSessions[cookieTok] = auth
  sesschange = true
  lsesslock.Unlock()
}

func findAuthSession(cookieTok string) *AuthSession{
  auth, ok := loginSessions[cookieTok]
  if !ok { return nil }
  if auth.Expires > 0 && auth.Expires < time.Now().Unix() {
    lsesslock.Lock()
    delete(loginSessions, cookieTok)
    sesschange = true
    lsesslock.Unlock()
    return nil
  }
  return auth
}

func deleteAuthSession(cookieTok string){
	lsesslock.Lock()
	delete(loginSessions, cookieTok)
	sesschange = true
	lsesslock.Unlock()
}

func writeAuthSessionCache() {
	if !sesschange { return } //no changes
	lsesslock.Lock()
	sesschange = false
	dat, err := json.Marshal(loginSessions)
	lsesslock.Unlock()
	if err == nil {
		err = ioutil.WriteFile(InstallDir+"/.session_cache", dat, 0600)
		if err != nil { fmt.Println("Cannot save session cache file!") }
	}
}

func readAuthSessionCache() {
	expire := time.Now().Add(time.Minute * 10).Unix()
	dat, err := ioutil.ReadFile(InstallDir+"/.session_cache")
	if err != nil { fmt.Println("Error reading session cache", err) ; return }
	_ = json.Unmarshal(dat, &loginSessions)
	for _, sess := range loginSessions {
		sess.Expires = expire
	}
	fmt.Println("Session Cache Length:", len(loginSessions) )
}
// ==============================
// Temporary Cookie Management
// ==============================
func GenerateClientCookieToken(w http.ResponseWriter) string {
	value := randomString(30)
	//fmt.Println("Set Cookie Value:", value)
	encoded, err := secure.Encode(cookiename, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  cookiename,
			Value: encoded,
			Path:  "/",
			Secure: true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
	return value
}

func ReadClientCookieToken(r *http.Request) string {
	cookie, err := r.Cookie(cookiename)
	if err != nil { return "" }
	var value string
	err = secure.Decode(cookiename, cookie.Value, &value)
	if err != nil { return "" }
	return value
}

func RemoveClientCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  cookiename,
		MaxAge: -1,	//This tells the browser to delete immediately
	}
	http.SetCookie(w, cookie)
}

// ==============================
//  Top-Level functions
// ==============================
func CheckLogin(w http.ResponseWriter, r *http.Request) (string, string) {
	//Return: userID, token of the login session
	userid := ""
	tok := ReadClientCookieToken(r)
	//fmt.Println("Check Login:", tok)
	
	if tok != "" {
		auth := findAuthSession(tok)
		if auth == nil { 
			//cookie expired or token de-authorized?
			RemoveClientCookie(w)
			tok = "" //expired tok - don't return this
		}else{
			userid = auth.Login_name
		}
	}
	if userid == "" {
		//No valid token/login yet - check basic auth parameters (but do not generate a token)
		uname, upass, ok := r.BasicAuth()
		if ok {
			userid = PerformLogin(uname, upass)
		}
	}
	return userid, tok
}

func GenerateLoginToken(w http.ResponseWriter, r *http.Request, user string) {
	//This is run after verifying login credentials with oauth, just to set the temporary token for this valid login
	tok := GenerateClientCookieToken(w)
	auth := AuthSession {
		Login_name: user,
	}
	saveAuthSession(tok, &auth)

}

func RemoveLoginToken(w http.ResponseWriter, r *http.Request, tok string) {
	deleteAuthSession(tok)
	RemoveClientCookie(w)
}
