package main

import (
	"log"

	"github.com/swkwon/gos/tcp"
)

func handler(c tcp.Context) {

}

func main() {
	server := tcp.New()
	server.RegisterHandler(handler)
	if e := server.Start("localhost:9999"); e != nil {
		log.Fatal(e)
	}
}
