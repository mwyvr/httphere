// httphere provides a simple http file server for ad-hoc use, testing HTML, etc.
//
// The package also provides an example for publishing to github, versioning, and configuration via the flag package.
//
// usage:
//   httphere
//
// httphere attempt to start the server listening on port 8080; if the port is
// already in use, it will iterate +1 port at a time until 8099.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	version string
)

func main() {
	var addr string

	fmt.Println(version)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	workdir, _ := os.Getwd()
	FileServer(r, "/", http.Dir(workdir))

	for i := 8080; i <= 8099; i++ {
		var err error
		addr = fmt.Sprintf("0.0.0.0:%d", i)
		fmt.Printf("Listening on: %s\n", addr)
		err = http.ListenAndServe(addr, r)
		if err != nil {
			fmt.Printf("Error, %s in use. Trying next.\n", addr) // most likely err
			continue
		}
		break
	}
	fmt.Println("No available ports between 8080 and 8099, quitting.")
}

// FileServer sets up a http.FileServer handler for an http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	// from the Chi examples:
	//	https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
