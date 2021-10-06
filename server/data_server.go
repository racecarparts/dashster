package server

import (
    "encoding/json"
    "fmt"
    "github.com/matryer/way"
    "github.com/racecarparts/dashster/service"
    "net/http"
)

type dataServer struct {
    Router *way.Router
}

func NewDataServer(router *way.Router) *dataServer {
    s := &dataServer{}
    s.Router = router
    s.dataServerRoutes()
    return s
}

func isErr(err error, errMsg string, httpStatusErr int, w http.ResponseWriter) bool {
   if err != nil {
       w.Write([]byte(fmt.Sprintf("{\"status\": %d,\"error\": \"%s\"}", httpStatusErr, errMsg)))
       return true
   }
   return false
}

func handleDataResp(jsonData []byte, err error, entityName string, w http.ResponseWriter) {
   if isErr(err, "unable to serialize " + entityName, http.StatusInternalServerError, w) {
       return
   }
   w.WriteHeader(http.StatusOK)
   w.Write(jsonData)
}

func (ds *dataServer) handleWorldClock() http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       times := service.WorldClock()
       timesData, err := json.Marshal(times)
       handleDataResp(timesData, err, "World Clock", w)
   }
}

func (ds *dataServer) handleCalendar() http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       dCal := service.Calendar()
       dCalJson, err := json.Marshal(dCal)
       handleDataResp(dCalJson, err, "Calendar", w)
   }
}

func (ds *dataServer) handleWeather() http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       weather := service.Weather()
       weatherJson, err := json.Marshal(weather)
       handleDataResp(weatherJson, err, "Weather", w)
   }
}

func (ds *dataServer) handleDocker() http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       docker := service.Docker()
       dockerJson, err := json.Marshal(docker)
       handleDataResp(dockerJson, err, "Docker", w)
   }
}

func (ds *dataServer) handlePullRequests() http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       prs := service.PullRequests()
       prJson, err := json.Marshal(prs)
       handleDataResp(prJson, err, "Pull Requests", w)
   }
}

func (ds *dataServer) handleMyCal() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        myCal := service.MyCal()
        myCalJson, err := json.Marshal(myCal)
        handleDataResp(myCalJson, err, "My Calendar", w)
    }
}
