package main

import (
	"compress/gzip"
	"log"

	"github.com/lesismal/nbio"
)

func main() {
	engine := nbio.NewEngine(nbio.Config{
		Network:            "tcp", //"udp", "unix"
		Addrs:              []string{":8888"},
		MaxWriteBufferSize: 6 * 1024 * 1024,
	})

	// hanlde new connection
	engine.OnOpen(func(c *nbio.Conn) {
		log.Println("OnOpen:", c.RemoteAddr().String())
	})
	// hanlde connection closed
	engine.OnClose(func(c *nbio.Conn, err error) {
		log.Println("OnClose:", c.RemoteAddr().String(), err)
	})
	// handle data
	engine.OnData(func(c *nbio.Conn, data []byte) {
		//c.Write(append([]byte{}, data...))
		header := []byte("HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nContent-Encoding: gzip\r\n\r\n")
		c.Write(header)

		xstr := []byte("<!DOCTYPE html><html><head></head><body><p>test</p></body></html>")
		w := gzip.NewWriter(c)
		w.Write(xstr)

		c.Close()
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
		return
	}
	defer engine.Stop()

	<-make(chan int)
}
