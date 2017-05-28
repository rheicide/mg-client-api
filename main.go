package main

import (
	"flag"
)

func main() {
	var addr string

	flag.StringVar(&addr, "addr", ":3000", "")
	flag.Parse()

	StartServer(NewServer(addr))
}
