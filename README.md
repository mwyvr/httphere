# httphere

`httphere` is a simple Go/Chi powered http server for ad hoc use such as
testing HTML or temporarily exposing a local file system.

## Installation

    go get -u https://github.com/solutionroute/httphere

## Usage

`httphere` will attempt to bind to a port between 8080 and 8099:

    httphere
