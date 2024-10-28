package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/axzilla/packify"
)

func main() {
	http.Handle("GET /", templ.Handler(packify.Index()))

	port := "8090"
	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(":"+port, nil)
}
