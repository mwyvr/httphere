/* httphere is a simple http file server for ad-hoc use such as testing HTML.

Given the intended use, cache-control headers are disabled by middleware by
default.

usage: httphere -h

httphere attempts to bind to the specified/default port; if the port can't be
bound to, it will attempt to bind to one of the next n ports.
*/

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	version     string = "1.0.2"
	httpAddress string = "0.0.0.0"
	httpPort    int    = 8080
	tryPorts    int    = 100
	noCache     bool   = true
)

func init() {
	flag.StringVar(&httpAddress, "address", httpAddress, "Address server should listen on")
	flag.IntVar(&httpPort, "port", httpPort, "Port server should bind to")
	flag.IntVar(&tryPorts, "tryPorts", tryPorts, "Attempt next n ports if specified can't be bound")
	flag.BoolVar(&noCache, "nocache", noCache, "Sets no-cache and other headers.\n-nocache=false to disable")
}

func main() {
	var addr string
	flag.Parse()

	r := chi.NewRouter()
	if noCache {
		r.Use(middleware.NoCache)
	}
	r.Use(middleware.Logger)

	// FileServer setup
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	FileServer(r, "/", http.Dir(workdir))

	// net.Listen to port or try range of ports if not able to bind
	for i := httpPort; i <= httpPort+tryPorts; i++ {
		addr = fmt.Sprintf("%s:%d", httpAddress, i)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Printf("Error: %s - trying next port.\n", err)
			continue
		}
		fmt.Printf("httphere (%s) listening at %s, serving files from %s\n", version, addr, workdir)
		err = http.Serve(l, r)
		if err != nil {
			panic(err)
		}
		break
	}
	fmt.Printf("Fatal: No bind-able ports between %d and %d.\n", httpPort, httpPort+tryPorts)
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
