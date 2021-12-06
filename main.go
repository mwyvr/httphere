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
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const nPorts int = 100

var (
	version     string = "development" // set via go build --ldflags, see Makefile
	httpAddress string = "0.0.0.0"
	httpPort    int    = 8080
)

func init() {
	flag.StringVar(&httpAddress, "address", httpAddress, "Address server should listen on")
	flag.IntVar(&httpPort, "port", httpPort, "Port server should bind to")
}

func main() {
	var addr string
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	FileServer(r, "/", http.Dir(workdir))

	for i := httpPort; i <= httpPort+nPorts; i++ {
		addr = fmt.Sprintf("%s:%d", httpAddress, i)
		fmt.Printf("httphere (%s) listening at: %s, serving files from: %s\n", version, addr, workdir)
		err = http.ListenAndServe(addr, r)
		if err != nil {
			fmt.Printf("Error: %s - trying next.\n", err)
			continue
		}
		break
	}
	fmt.Printf("No available ports between %d and %d.", httpPort, httpPort+nPorts)
}

// FileServer sets up a http.FileServer handler for an http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	// from the Chi examples:
	//	https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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
