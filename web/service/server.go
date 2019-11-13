package service

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
		Directory:  "templates",
		Extensions: []string{".html"},
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
		}
	}

	mx.HandleFunc("/", templateHandler(formatter)).Methods("GET")
	mx.HandleFunc("/unknown", notImplemented).Methods("GET")
	mx.HandleFunc("/api/test", jsonHandler(formatter)).Methods("GET")
	mx.HandleFunc("/submit", submitHandler(formatter)).Methods("POST")
	mx.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir(webRoot+"/assets/"))))

}

func notImplemented(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "501 Not Implemented\nWe are still working on it. Please wait!", 501)
}
