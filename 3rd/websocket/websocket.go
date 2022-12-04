package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "httprouter service address")

func BuildWebsocketClient() {
	flag.Parse()
	log.SetFlags(0)
	irq := make(chan os.Signal, 1)
	signal.Notify(irq, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial websocket occurs error:", err)
	}
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("reading websocket occurs error:", err)
			}
			log.Printf("websocket recv: %s", message)
		}
	}()
}

func main() {
	BuildWebsocketClient()
}
