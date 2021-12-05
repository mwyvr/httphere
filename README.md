# httphere

`httphere` is a simple Go/Chi powered http server for ad hoc use such as
testing HTML or temporarily exposing a local file system at the current
working directory (hence http**here**) and below.

## Installation

    go get -u https://github.com/solutionroute/httphere

## Usage

`httphere` with no flags will attempt to start a webserver in the current working 
directory, binding to the default or specified port at all available IPv4 addresses 
on the machine including `localhost` via the special address `0.0.0.0`.  

Available flags:

    -address string
            Address server should listen on (default "0.0.0.0")
    -port int
            Port server should bind to (default 8080)

            If the port can't be bound to (e.g. if already a bound to a server), 
            ports n+100 will be tried.
