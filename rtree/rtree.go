package main

import (
	"flag"
	"fmt"
	"github.com/fumin/rtree"
)

const (
	defaultPort = 6389
	portUsage   = "port to bind to"
)

var port int

func init() {
	flag.IntVar(&port, "port", defaultPort, portUsage)
	flag.IntVar(&port, "p", defaultPort, portUsage)
}

func main() {
	flag.Parse()
	if port == defaultPort {
		fmt.Println("Binding to default port", defaultPort)
	}

	s, err := rtree.NewServer("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	s.LoopAccept()
}
