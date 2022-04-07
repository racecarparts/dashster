package server

import "net/http"

func (s *viewServer) fileServerRoutes() {
	s.Router.HandleFunc("GET", "/view...", s.handleView())
	s.Router.Handle("GET", "/static...", http.FileServer(s.handleStatic()))
}

func (ds *dataServer) dataServerRoutes() {
	ds.Router.HandleFunc("GET", "/worldclock", ds.handleWorldClock())
	ds.Router.HandleFunc("GET", "/calendar", ds.handleCalendar())
	ds.Router.HandleFunc("GET", "/weather", ds.handleWeather())
	ds.Router.HandleFunc("GET", "/docker", ds.handleDocker())
	ds.Router.HandleFunc("GET", "/pullrequests", ds.handlePullRequests())
	ds.Router.HandleFunc("GET", "/mycal", ds.handleMyCal())
}
