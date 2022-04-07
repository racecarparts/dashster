package server

import (
	"embed"
	"fmt"
	"github.com/matryer/way"
	htmlTemplate "html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	viewDir   = "public/view/*"
	staticDir = "public"
)

var (
	PublicFiles   embed.FS
	initTemplates sync.Once
	initStatic    sync.Once
	viewTemplates *htmlTemplate.Template
)

type viewServer struct {
	Router    *way.Router
	embedRoot string
}

func NewViewServer(router *way.Router) *viewServer {
	s := &viewServer{}
	s.Router = router
	s.fileServerRoutes()
	return s
}

func (s *viewServer) handleStatic() http.FileSystem {
	staticFiles, err := fs.Sub(PublicFiles, staticDir)
	if err != nil {
		log.Printf("static files not found in %s", staticDir)
		return nil
	}
	return http.FS(staticFiles)
}

func (s *viewServer) handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTemplates.Do(func() {
			var err error
			viewTemplates, err = htmlTemplate.ParseFS(PublicFiles, viewDir)
			if err != nil {
				log.Printf("html templates not found in %s", viewDir)
				return
			}
			err = viewTemplates.Execute(w, map[string]string{})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})

		tplName := strings.SplitAfter(r.URL.Path, "/view/")
		tpl := viewTemplates.Lookup(tplName[1])
		fmt.Println(r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	}
}
