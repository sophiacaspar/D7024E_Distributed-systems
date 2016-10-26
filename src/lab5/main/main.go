package main

import (
	dht "lab5/dht"
	"os"
)

func main() {
	port := os.Args[1]
	dht.StartUpNetwork(port)

}