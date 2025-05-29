package main

import (
	"flag"
)

var MODE bool //isServer

func init() {
	// flag.BoolVar(&MODE, "server", true, "run as server")
	flag.BoolVar(&MODE, "S", true, "run as server")
	// flag.BoolVar(&MODE, "client", false, "run as client")
	flag.BoolVar(&MODE, "C", false, "run as client")
}

func main() {
	flag.Parse()
	if !MODE {
		runServer()
	} else {
		runClient()
	}
}
