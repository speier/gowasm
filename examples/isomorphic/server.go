package main

//go:generate rice embed-go

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/speier/gowasm/examples/isomorphic/app"
	"github.com/speier/gowasm/pkg/server"
)

var (
	addr = flag.String("addr", ":8080", "listen address")
)

func init() {
	flag.Parse()
}

func main() {
	box := rice.MustFindBox("static")

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.URLFormat)

	r.Handle("/static/*", http.StripPrefix("/static/", static(box)))
	r.NotFound(index(box))

	log.Printf("Server listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

func index(box *rice.Box) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initState := &app.State{Count: 1}
		view := app.View(initState, nil)
		html := server.RenderToString(view)
		tmpl := template.Must(template.New("").Parse(box.MustString("index.html")))
		jsonInitState, err := json.Marshal(initState)
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, map[string]interface{}{
			"appMarkup": html,
			"initState": string(jsonInitState),
		})
	})
}

func static(box *rice.Box) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(box.HTTPBox()).ServeHTTP(w, r)
	})
}
