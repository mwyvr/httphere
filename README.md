# httphere

`httphere` is a simple Go/Chi powered http file server intended for ad hoc use
such as testing HTML or temporarily exposing a local file system at the current
working directory (hence http**here**) and below.

Given the expected use, cache-control headers are disabled by default.

## Installation

    go install github.com/solutionroute/httphere@latest

## Usage

`httphere` starts a webserver in the current working directory at the default
or specified port; if that port cannot be bound to, it will attempt to bind to
one of the next 100 ports in sequence.

By default, via the special address `0.0.0.0`, the server binds to all available
IPv4 addresses on the machine including its real IP address(es) and localhost/loopback
at 127.0.0.1.

Available flags:

    -address string
            Address server should listen on (default "0.0.0.0")
    -nocache
            Sets no-cache and other headers.
            -nocache=false to disable (default true)
    -port int
            Port server should bind to (default 8080)
