# Generic template for a go-app based WASM website
This is a generic template that can be copied/adjusted to get a WASM-based website up and running quickly. This already has the authentication/token management systems in place for HTTP/REST API calls, it just needs the actual username/password login handling filled in based on the site requirements.

## Build setup
There is a small "Makefile" in this directory to simplify building/deploying the site. If "make" cannot be used on this system,  please look inside that file to see the matching "go run" commands for each step.

* `make`
   * Build the site binary and put it into the current directory
* `make package`
   * Run this after `make`
   * Assembles a subdirectory called "dist" which contains everything needed to run the site.
* `make install`
   * Run this after `make package`, typically as the root user on the system (`sudo make install`)
   * This copies the "dist" directory to "/usr/share/<name>", and symlinks the binary to "/usr/bin/<name>" so it can be easily started.
   * NOTE: If this is an update, it will completely remove/replace any older installed version of the utility. Make sure to "restart" the binary if this is 
* `make clean`
   * Remove all compiled binaries and package directories - returning the directory back to it's pristine context.

### Typical routine/script for updating a running service
```
#!/bin/sh
cd <source directory>
make
if [ $? -eq 0 ] ; then
  make package
  <stop running service>
  sudo make install
  <restart service>
fi
```

This will result in a minimal amount of downtime, as the time it takes to copy the files over is all that is needed.
It is possible that the "stop service" step can be omitted in some cases where there are very few files in the dist structure, and a simple restart is all that is needed.

## Layout
* **src** directory
   * This contains all the server-side code for the website. 
   * This code is PRIVATE, and will not be exposed to the user/browser.
* **src-wasm** directory
   * This contains all the go-code used to generate/build the WASM binary used to display the website.
   * This code gets compiled into a web assembly binary. While abstracted, this *is* viewable and changeable by the user/browser.
* **src-common** directory
   * This directory contains code common to both the server and the client/browser.
   * This is very useful for defining things like API structure formats and such so that you don't need to maintain two different copies of the structures.
   * All "*.go" files within this directory are sym-linked into both the src/ and src-wasm/ directories before the build is started.
   * Running `make clean` will run the cleanup routine and remove the symlinks from both directories.
* **web** directory
   * This contains all the static files served via the site
   * "style.css" : This is the CSS file used for stying the site. Most styling should be done in the go-app sources directly, but this is good for often-used styling classes and color themes and such.
   * "favicon.ico" : The favicon for the site. This *can* be a PNG file instead of an ICO file in my experience.
   * "static" directory : Everything in this directory is publicly available via the "<site:port>/static/<path/filename>" URL. This is useful for self-hosted static files like images/icons/etc.


## Initial Setup
* **build.go**
   * Adjust the "pName" global string at the top of the file to the name of your program/site binary. This will be used in many places for settings files and such.
* **src/main.go**
   * Adjust the "progname" global string at the top of the file to the name of your program/site binary. This should match the binary name changed in build.go 
   * Optionally adjust the app.Handler definitions (Title/Author/Name/ShortName). These are seen by the user/browser.
* **src/auth_session.go**
   * The "PerformLogin()" function at the top of the file needs to be filled in as needed. That is where you add in the authentication backend usage as needed.
* **src/api.go**
   * The "Evaluate()" function at the top of the file needs to be filled in the specific API calls you want to support.
   * The "login" and "logout" API's are already built-in to the authentication systems.
   * Note: All the APIs are accessed via `http://<site_location:port>/api/<api_name>`


## OAuth Support
For OAuth login support, there needs to be a couple more changes, detailed here:

1. **src/auth_oauth.go**
   * Fill in the "StartOauthLogin()" function near the top of the file which sets up and issues the redirect for the designated oauth provider.
   * Fill in the "HandleOauthLogin()" function near the top of the file which reads back the info from the OAuth provider and generates a new session token for the user if the login was accepted by the OAuth system.
