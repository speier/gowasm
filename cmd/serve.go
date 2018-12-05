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
	reload bool
)

func init() {
	serveCmd.PersistentFlags().StringVarP(&addr, "listen", "l", ":8010", "listen address")
	serveCmd.PersistentFlags().BoolVarP(&reload, "reload", "r", true, "reload changes")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve [directory]",
	Short: "Run server for wasm",

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}

		dir, err := filepath.Abs(dir)
		if err != nil {
			log.Fatal("invalid directory", err)
		}

		name := tempFilename("gowasm_", ".wasm")
		defer os.Remove(name)

		box := rice.MustFindBox("static")
		r := chi.NewRouter()

		r.Handle("/app.wasm", wasm(dir, name))
		r.Handle("/static/*", http.StripPrefix("/static/", static(box)))
		r.NotFound(index(dir, box))

		httpServer := &http.Server{Addr: addr, Handler: r}
		defer shutdown(httpServer)
		go start(dir, httpServer)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
	},
}

func start(dir string, httpServer *http.Server) {
	a := httpServer.Addr
	if strings.HasPrefix(a, ":") {
		a = "localhost" + a
	}
	log.Printf("Server listening on http://%s", a)
	log.Printf("Serving directory %s", dir)
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

func index(dir string, box *rice.Box) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// project have it's own index.html
		projhtml := projectHtml(dir)
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

func wasm(dir string, name string) http.HandlerFunc {
	build(dir, name, true)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		build(dir, name, false)
		http.ServeFile(w, r, name)
	})
}

func build(dir string, name string, exit bool) {
	log.Println("Building...")
	start := time.Now()
	res, err := run(fmt.Sprintf("go build -o %s %s", name, dir), "GOOS=js", "GOARCH=wasm")
	elapsed := time.Since(start)
	if err != nil {
		log.Println("Build failed", res)
		if exit {
			os.Exit(1)
		}
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

func projectHtml(dir string) string {
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
