# dashster

A slapped together dashboard app for seeing stuff that matters to a developer.

This is a personal project that I mess with from time to time.  Depending on my 
employer and what git service is used, I've stitched some things together.  
If anything, it's a good jumping off point if you want to build a simple app that 
can call web APIs, command line utils, etc and display the results all in one place. 

Written in Go, HTML, JS, and CSS.  Uses WebView for rendering the UI.  I use it on 
macOS, but it should work on any platform that Go supports.

## Build it...

### As a macOS App:
```
$ ./build-macos.sh
```

### For anything else:
```
$ go build
```

## Run it...
On first run, the app will create a config file in your home directory:
* macOS: `~/.dashster_config.json`
* linux: `~/.dashster_config.json`
* windows: 
    * `c:\Users\[username]~/.dashster_config.json`
    * `~/.dashster_config.json`

Then update the config file with what you need. 

## Development
I'm not offering much help here right now.  Maybe someday I'll clean this up, but for now it's just a rough personal project. I chose to use plain JS, HTML, CSS, Bootstrap, and Go templates. Why? It's a local desktop app. It works. Moving on.

There are several modules, some are commented out because 
I don't use them or I changed my mind.

Some of the modules (some are disabled in code):
* 3 month calendar
* Local machine docker processes
* Gitlab MRs
* Github PRs (disabled)
* Personal calendar agenda list (macOS specific - disabled)
* Weather
* World clock with daylight tracker