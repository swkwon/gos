package main

import (
	"gos/tcp"
	"log"
)

func handler(s *tcp.Session, data []byte) {

}

func main() {
	server := tcp.New()
	server.RegisterHandler(handler)
	if e := server.Start(":9999"); e != nil {
		log.Fatal(e)
	}
}
