package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/axzilla/goilerplate/assets"
	"github.com/axzilla/packify/pages"
)

func main() {
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)

	mux.Handle("GET /", templ.Handler(pages.Index()))

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
