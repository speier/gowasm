// +build !wasm

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/speier/gowasm/pkg/server"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	http.ServeFile(w, r, "main.wasm")
}

func ssrHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	initState := &State{Count: 0}
	app := App(initState)
	html := server.RenderToString(app)

	fmt.Fprintf(w, html)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/ssr", ssrHandler)
	mux.HandleFunc("/main.wasm", wasmHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
