package main

import (
	"bytes"
	"fmt"
	"strings"

	// "io"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/axzilla/goilerplate/assets"
	"github.com/axzilla/packify/pages"
	"github.com/axzilla/packify/utils"
)

func main() {
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)

	mux.Handle("GET /", templ.Handler(pages.Index("")))
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		include := "*"
		repoUrl := strings.ReplaceAll(r.FormValue("url"), " ", "")

		if !utils.IsValidGithubURL(repoUrl) {
			err := fmt.Errorf("Not a valid GitHub repo URL")
			pages.Index(err.Error()).Render(r.Context(), w)
			return
		}

		if r.FormValue("include") != "" {
			include = strings.ReplaceAll(r.FormValue("include"), " ", "")
		}
		exclude := strings.ReplaceAll(r.FormValue("exclude"), " ", "")

		includePatterns := strings.Split(include, ",")
		excludePatterns := strings.Split(exclude, ",")

		var filetreeBuffer bytes.Buffer
		var contentsBuffer bytes.Buffer

		fileSystem, err := utils.FileSystem(repoUrl)
		if err != nil {
			pages.Index(err.Error()).Render(r.Context(), w)
			return
		}

		err = utils.WriteToBuffer(fileSystem, &repoUrl, includePatterns, excludePatterns, &filetreeBuffer, &contentsBuffer)
		if err != nil {
			pages.Index(err.Error()).Render(r.Context(), w)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=packify.txt")

		filetreeBuffer.WriteTo(w)
		contentsBuffer.WriteTo(w)
	})

	port := "8090"
	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(":"+port, mux)
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"
	// We need this for Templ to work
	disableCacheInDevMode := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isDevelopment {
				w.Header().Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, r)
		})
	}
	// Serve static files from the assets directory
	var fs http.Handler
	if isDevelopment {
		fs = http.FileServer(http.Dir("./assets"))
	} else {
		fs = http.FileServer(http.FS(assets.Assets))
	}
	mux.Handle("GET /assets/*", disableCacheInDevMode(http.StripPrefix("/assets/", fs)))
}
