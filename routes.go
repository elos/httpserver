package httpserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"go/build"

	"github.com/elos/data"
	"github.com/elos/models"
	t "github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

var (
	assetsDir    = filepath.Join(defaultBase("github.com/elos/httpserver"), "assets")
	templatesDir = filepath.Join(assetsDir, "templates")
	imgDir       = filepath.Join(assetsDir, "img")
	cssDir       = filepath.Join(assetsDir, "css")
)

func setupRoutes(s *HTTPServer) {
	s.GET("/", Template("index.html"))
	s.GET("/sign-in", Template("sign-in.html"))
	s.POST("/sign-in", Auth(SignInHandle, t.Auth(t.FormCredentialer), s.Store))

	s.ServeFiles("/css/*filepath", http.Dir(cssDir))
	s.ServeFiles("/img/*filepath", http.Dir(imgDir))
}

func defaultBase(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

type Page struct {
}

func Template(name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := filepath.Join(templatesDir, name)
		t, err := template.ParseFiles(path)
		if err != nil {
			log.Print("Template error: %s", err)
			return
		}
		t.Execute(w, Page{})
	}
}

func SignInHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	w.Write([]byte("hello" + a.Client().(models.User).Name()))
}
