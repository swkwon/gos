package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/swkwon/gos/tcp"
)

func handler(c tcp.Context) {
	d := c.GetReceived()
	log.Println(string(d))
	m := make(map[string]interface{})
	json.Unmarshal(d, &m)
	c.JSON(&m)
}

func main() {
	go func() {
		log.Println("start client")
		time.Sleep(1 * time.Second)
		c, e := net.Dial("tcp", "localhost:9999")
		if e != nil {
			log.Println(e)
			return
		}
		msg := `{"message":"hello world"}`
		header := make([]byte, 4)
		binary.LittleEndian.PutUint32(header, uint32(len(msg)))
		b := []byte(msg)
		packed := append(header, b...)
		c.Write(packed)
		buffer := make([]byte, 1024)
		c.Read(buffer)
		size := binary.LittleEndian.Uint32(buffer[:4])
		received := buffer[4 : 4+size]
		log.Println("received", string(received))
		log.Println("end client")
	}()

	server := tcp.New()
	server.RegisterHandler(handler)
	if e := server.Start("localhost:9999"); e != nil {
		log.Fatal(e)
	}
}
