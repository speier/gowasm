package cmd

//go:generate rice embed-go

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

var (
	addr   string
	dir    string
	reload bool
)

func init() {
	serveCmd.PersistentFlags().StringVarP(&addr, "addr", "a", ":8010", "listen address")
	serveCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "directory to serve")
	serveCmd.PersistentFlags().BoolVarP(&reload, "reload", "r", true, "reload changes")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server for wasm",

	Run: func(cmd *cobra.Command, args []string) {
		box := rice.MustFindBox("static")

		r := chi.NewRouter()

		name := tempFilename("gowasm_", ".wasm")
		defer os.Remove(name)

		r.Handle("/app.wasm", wasm(name))
		r.Handle("/static/*", http.StripPrefix("/static/", static(box)))
		r.NotFound(index(box))

		httpServer := &http.Server{Addr: addr, Handler: r}
		defer shutdown(httpServer)
		go start(httpServer)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
	},
}

func start(httpServer *http.Server) {
	log.Printf("Server listening on %s", addr)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("http server failed to serve", err)
	}
}

func shutdown(httpServer *http.Server) {
	log.Println("Server shutdown")
	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal("error happened during http server shutdown", err)
	}
}

func index(box *rice.Box) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// project have it's own index.html
		projhtml := projectHtml()
		if len(projhtml) > 0 {
			fmt.Fprint(w, projhtml)
		} else {
			// use default builtin index.html for dev
			tmpl := template.Must(template.New("").Parse(box.MustString("index.html")))
			tmpl.Execute(w, map[string]interface{}{
				"title": "GOWASM CLI",
			})
		}
	})
}

func static(box *rice.Box) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(box.HTTPBox()).ServeHTTP(w, r)
	})
}

func wasm(name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		build(name)
		http.ServeFile(w, r, name)
	})
}

func build(name string) {
	log.Println("Building...")
	start := time.Now()
	res, err := run(fmt.Sprintf("go build -o %s %s", name, dir), "GOOS=js", "GOARCH=wasm")
	elapsed := time.Since(start)
	if err != nil {
		log.Println("Build failed", res, err.Error())
	} else {
		log.Println("Build done in", elapsed)
	}
}

func run(command string, env ...string) (string, error) {
	prts := strings.Fields(command)
	name := prts[0]
	args := prts[1:]

	cmd := exec.Command(name, args...)
	cmd.Dir, _ = os.Getwd()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	res, err := cmd.CombinedOutput()
	return string(res), err
}

func tempFilename(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func projectHtml() string {
	name := filepath.Join(dir, "index.html") // other dirs/files and/or configurable?

	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return ""
	}

	f, err := os.Open(name)
	if err != nil {
		return ""
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return ""
	}

	inject := `<script src="/static/wasm_exec.js"></script><script src="/static/wasm_inst.js"></script>`
	doc.Find("body").AppendHtml(inject)

	ret, err := doc.Html()
	if err != nil {
		return ""
	}

	return ret
}
