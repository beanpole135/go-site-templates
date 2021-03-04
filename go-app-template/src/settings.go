package main

import (
	"fmt"
	"encoding/json"
	"os"
	"path/filepath"
	"io/ioutil"
)

var InstallDir string = getInstallDir()

//Settings structure, only used internally by Go
type Settings struct {
	FilePath string			`json:"-"`	//internal - just for keeping track of where the settings were read from
	ListenURL string		`json:"listen_url"`
	LSessHashKey []byte	`json:"login_session_hash_key"`
	LSessBlockKey []byte	`json:"login_session_block_key"`
}

func (S *Settings) LoadDefaults() {
	if S.ListenURL == "" { S.ListenURL = "127.0.0.1:8080" }
}

func readSettings(path string) Settings {
  var settings Settings
  if(path == "" ){
    check := []string{ InstallDir+"/config.json", "/usr/etc/"+progname+"/config.json", "/etc/"+progname+"/config.json" }
    for _, pathchk := range(check) {
      if fileExists(pathchk) { path = pathchk; break }
      fmt.Println("Path did not exist:", pathchk)
    }
  }
  if path != "" {
    fmt.Println("Using Settings:", path)
    dat, err := ioutil.ReadFile(path)
    if err == nil {
      err = json.Unmarshal(dat, &settings)
      if err != nil { fmt.Println("Could not parse JSON config!", err) ; os.Exit(1) }
	settings.FilePath = path
    }else{
      fmt.Println("Could not read settings:", err)
    }
  } else {
    fmt.Println("No config file specified")
  }
  // Verify all the settings
  settings.LoadDefaults()
  return settings
}

func getInstallDir() string {
    installdir, _ := os.Executable()
    installdir, _ = filepath.Abs(installdir)
    installdir, _ = filepath.EvalSymlinks(installdir)
    installdir = filepath.Dir(installdir)
  fmt.Println("Loading files from directory:", installdir)
  return installdir
}

func (S *Settings) Save() {
	if S.FilePath == "" { return } //no config file to save
	fmt.Println("Updating config file:", S.FilePath)
	dat, err := json.MarshalIndent(S, "", "  ")
	if err == nil {
		err = os.Rename(S.FilePath, S.FilePath+".prev")
		if err != nil { fmt.Println("Cannot rename old config file - stopping save:", err) ; return }
		err = ioutil.WriteFile(S.FilePath, dat, 0600)
		if err != nil { 
			fmt.Println("Cannot write config file - restoring previous:", err)
			os.Rename(S.FilePath+".prev", S.FilePath)
		}
	}else {
		fmt.Println(" - error marshalling json settings:", err)
	}
}
