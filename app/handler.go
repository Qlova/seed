package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/qlova/seed/script"
)

var intranet, _ = regexp.Compile(`(^192\.168\.([0-9]|[0-9][0-9]|[0-2][0-5][0-5])\.([0-9]|[0-9][0-9]|[0-2][0-5][0-5]):.*$)`)

func isLocal(r *http.Request) (local bool) {
	local = strings.Contains(r.RemoteAddr, "[::1]") || strings.Contains(r.RemoteAddr, "127.0.0.1")
	if intranet.Match([]byte(r.RemoteAddr)) {
		local = true
	}

	split := strings.Split(r.Host, ":")
	if len(split) == 0 {
		local = false
	} else {
		if split[0] != "localhost" {
			local = false
		}
	}

	return
}

//Handler returns an http.Handler that serve's the app.
func (app App) Handler() http.Handler {
	router := http.NewServeMux()

	app.build()

	var rendered = app.document.Render()

	var document, err = mini(rendered)
	if err != nil {
		document = rendered
	}

	var scripts = script.Scripts(app.document)

	var worker = app.worker.Render()

	router.Handle("/Qlovaseed.png", gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		icon, _ := fsByte(false, "/Qlovaseed.png")
		w.Write(icon)
		return
	})))

	router.Handle("/call/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		script.Handler(w, r, r.URL.Path[6:])
	}))

	router.Handle("/seed.socket", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isLocal(r) {
			localClients++
			singleLocalConnection = localClients == 1
			socket(w, r)
		}
	}))

	var manifest = app.manifest.Render()
	router.Handle("/app.webmanifest", gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Write(manifest)
	})))

	router.Handle("/index.js", gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isLocal(r) {
			//Don't use a web worker if we are running locally.
			w.Header().Set("content-type", "text/javascript")
			w.Write([]byte(`self.addEventListener('install', () => {self.skipWaiting();});`))
		} else {
			w.Header().Set("content-type", "text/javascript")
			w.Write(worker)
		}
	})))

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if content, ok := scripts[r.URL.Path]; ok {
			w.Header().Set("Content-Type", "application/js")
			fmt.Fprint(w, content)
			return
		}

		w.Write(document)
	}))

	return router
}