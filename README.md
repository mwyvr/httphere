# httphere

`httphere` is a simple Go/Chi powered http file server for ad hoc use such as
testing HTML or temporarily exposing a local file system at the current
working directory (hence http**here**) and below.

## Installation

    go install github.com/solutionroute/httphere@latest

## Usage

`httphere` with no flags starts a webserver in the current working directory.
By default, via the special address `0.0.0.0`, the server binds to all available
IPv4 addresses on the machine.

Available flags:

    -address string
            Address server should listen on (default "0.0.0.0")
    -port int
            Port server should bind to (default 8080)

            If the port can't be bound to, a range of 100 incremental
            ports will be tried.
