package main

import (
	"embed"
	"fmt"
	"github.com/matryer/way"
	"github.com/racecarparts/dashster/server"
	"github.com/racecarparts/dashster/service"
	"github.com/webview/webview"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	//go:embed public
	EmbeddedFiles embed.FS
	RepoRoot string
	events chan string // keyboard events
)

func init() {
	// events is a channel of string events that come from the front end
	events = make(chan string, 1000)
	RepoRoot = "public"
	server.PublicFiles = EmbeddedFiles
}

func run() error {
	err := service.ReadOrCreateConfig()
	if err != nil {
		return err
	}
	err = setupUI(events)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func setupUI(events chan string) error {
	// channel to get the web prefix
	prefixChannel := make(chan string)
	// run the web server in a separate goroutine
	go app(prefixChannel)
	prefix := <-prefixChannel

	debug := true
	w := webview.New(debug)
	w = service.RegisterBindings(w)

	defer w.Destroy()
	w.SetTitle("Dashster")
	w.SetSize(1000, 1600, webview.HintNone)
	w.Navigate(prefix + "/view/index")
	fmt.Println(prefix + "/view/index")
	w.Run()

	return nil
}

func app(prefixChannel chan string) {
	router := way.NewRouter()
	server.NewViewServer(router)
	server.NewDataServer(router)


	http.FileServer(http.FS(server.PublicFiles))

	// get an ephemeral port, so we're guaranteed not to conflict with anything else
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	portAddress := listener.Addr().String()
	prefixChannel <- "http://" + portAddress
	listener.Close()
	server := &http.Server{
		Addr:    portAddress,
		Handler: router,
	}
	server.WriteTimeout = 3 * time.Minute
	server.ReadTimeout = 3 * time.Minute
	server.ListenAndServe()
}


