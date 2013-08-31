package main

import (
  "flag"
  "fmt"
  "github.com/fumin/rtree"
)

var port int

func init() {
  const (
    defaultPort = 6389
    portUsage = "port to bind to"
  )
  flag.IntVar(&port, "port", defaultPort, portUsage)
  flag.IntVar(&port, "p", defaultPort, portUsage)
}

func main() {
  flag.Parse()
  s, err := rtree.NewServer("tcp", fmt.Sprintf(":%d", port))
  if err != nil { panic(err) }
  s.LoopAccept()
}
