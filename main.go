/* httphere is a simple http file server for ad-hoc use such as testing HTML.

httphere attempts to bind to the specified/default port; if the port can't be
bound to, it will attempt to bind to one of the next 99 ports. no-cache headers
are set by default.

httphere -h for usage
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	version     string = "1.1.1"
	httpAddress string = "0.0.0.0"
	httpPort    int    = 8080
	tryPorts    int    = 100
	noCache     bool   = true
)

func init() {
	flag.StringVar(&httpAddress, "address", httpAddress, "Address server listens on")
	flag.IntVar(&httpPort, "port", httpPort, "Port server should bind to")
	flag.BoolVar(&noCache, "nocache", noCache, "Set no-cache headers.\n-nocache=false to disable")
}

func main() {
	var addr string
	flag.Parse()

	r := chi.NewRouter()
	if noCache {
		r.Use(middleware.NoCache)
	}
	r.Use(middleware.Logger)

	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	FileServer(r, "/", http.Dir(workdir))

	// bind to specified port or to port+n if unavailable
	for i := httpPort; i <= httpPort+tryPorts-1; i++ {
		addr = fmt.Sprintf("%s:%d", httpAddress, i)
		log.Printf("attempting to bind to %s\n", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("%s, trying next port.\n", err)
			continue
		}
		log.Printf("httphere v%s listening at %s, serving files from %s\n", version, addr, workdir)
		err = http.Serve(listener, r)
		if err != nil {
			log.Fatalf("fatal: %s", err)
		}
		break
	}
	log.Fatalf("fatal: No bind-able ports between %d and %d.\n", httpPort, httpPort+tryPorts)
}

// FileServer sets up a http.FileServer handler for an http.FileSystem.
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Fatal("FileServer does not permit URL parameters")
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
