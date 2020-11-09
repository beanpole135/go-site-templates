// +build ignore
package main

import (
	"os"
	"os/exec"
	"fmt"
	"io/ioutil"
	"strings"
)

//Product name (static)
var pName string =  "CHANGEME"


func exit_err(err error, details string){
  if(err != nil){
    fmt.Println("[ERROR] ", details);
    fmt.Println("  ", err)
    os.Exit(1)
  }
}

func runCMD( command []string, workdir string, envs ...string) error {
  wd, _ := os.Getwd() //get working directory
  cmd := exec.Command( command[0], command[1:]...)
  if workdir != "" { cmd.Dir = workdir }
  cmd.Env = append(os.Environ(), "GOPATH="+os.Getenv("GOPATH")+":"+wd+workdir)
  for _, env := range envs { cmd.Env = append(cmd.Env, env) }
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Run()
  return err
}

func symlinkAll(fromdir string, todir string){
  wd, _ := os.Getwd() //get working directory
  files, _ := ioutil.ReadDir(fromdir)
  for _, file := range files {
    if !strings.HasSuffix(file.Name(), ".go") { continue }
    fromfile := wd+"/"+fromdir+"/"+file.Name()
    tofile := todir+"/"+file.Name()
    if _, err := os.Stat(tofile) ; !os.IsNotExist(err) { return }
    err := os.Symlink(fromfile, tofile)
    if err != nil { fmt.Println("Symlink Error:", err) }
  }
}

func localCopy(from string, to string){
  if _, err := os.Stat(to) ; !os.IsExist(err) {
    os.RemoveAll(to)
  }
  cmd := exec.Command("cp", "-rf", from, to)
  err := cmd.Run()
  exit_err(err, "Could not copy: "+ from+ " -to- "+to)
}

func doClean(){
  files := []string{ pName, "web/app.wasm" } //files to remove
  dirs := []string{ "dist" } //dirs to remove
  linkdirs := []string{ "src", "src-wasm"}
  //Remove files
  for _, file := range files {
    if _, err := os.Stat(file) ; !os.IsExist(err) { os.Remove(file) }
  }
  //Remove dirs
  for _, dir := range dirs {
    if _, err := os.Stat(dir) ; !os.IsExist(err) { os.RemoveAll(dir) }
  }
  //Remove symlinks in dirs
  for _, dir := range linkdirs {
    flist, _ := ioutil.ReadDir(dir)
    for _, link := range flist {
      if (link.Mode() & os.ModeSymlink != 0) { os.Remove(dir+"/"+link.Name()) }
    }
  }
}

func doInstall() error {
  if _, err := os.Stat("dist") ; !os.IsExist(err) {
    exit_err(err, "Please run \"make package\" first")
  }
  installdir := "/usr/share/"+pName
  os.RemoveAll(installdir)
  localCopy("dist", installdir)
  //Now make the symlink for the binary
  if _, err := os.Stat("/usr/bin/"+pName) ; os.IsNotExist(err) {
    err = os.Symlink(installdir+"/"+pName, "/usr/bin/"+pName)
    exit_err(err, "Could not symlink the binary into /usr/bin")
  }
  return nil
}

func doPackage() error{
  err := os.MkdirAll("dist", 0755)
  exit_err(err, "Could not create \"dist\" directory")
  localCopy(pName, "dist/"+pName)
  localCopy("web", "dist/web")
  return err
}

func buildWasm() error {
  fmt.Println(" - building WASM")
  symlinkAll("src-common", "src-wasm")
  err := runCMD( []string{"go","get"}, "src-wasm")
  //exit_err(err, "WASM go get failed")
  err = runCMD( []string{"go","build","-o","app.wasm"}, "src-wasm", "GOOS=js", "GOARCH=wasm")
  err = os.Rename("src-wasm/app.wasm", "web/app.wasm")
  return err
}

func doBuild() error {
  fmt.Println(" - building server")
  symlinkAll("src-common", "src")
  err := runCMD( []string{"go","get"}, "src")
  //exit_err(err, "go get failed")
  berr := runCMD( []string{"go","build","-o",pName}, "src")
  if _, terr := os.Stat("src/"+pName) ; !os.IsExist(terr) {
    err = os.Rename("src/"+pName, pName)
  }else{
    exit_err(berr, "Build Failed")
  }
  if err == nil { err = buildWasm() }
  return err
}

// This is the build routine for each of the projects
func main(){
  subcmd := "build"
  if(len(os.Args)>=2){ subcmd = os.Args[1] }
  fmt.Println(subcmd+"ing "+pName+"...");
  //Now run the appropriate type of operation
  var err error
  err = nil
  switch(subcmd){
    case "build":
        err = doBuild()

    case "clean":
        doClean()

    case "package":
        err = doPackage()

    case "install":
        err = doInstall()

    default:
        fmt.Println("Unknown action: ", subcmd)
	fmt.Println("Available actions are:")
	fmt.Println(" - make build:", "Compile the tools for the current system OS/ARCH")
	fmt.Println(" - make clean:", "Cleanup all the build files")
	fmt.Println(" - make package:", "Create a sterile \"dist\" directory ready to be copied/installed someplace")
	fmt.Println(" - make install:", "Install the package output to the designated directory")
      os.Exit(1)
  }
  if(err != nil){ 
    fmt.Println("[Error]", err)
    os.Exit(1) 
  } else {
    fmt.Println("[Success]")
    os.Exit(0)
  }
}
