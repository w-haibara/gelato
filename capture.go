package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"image"
	"image/jpeg"

	"github.com/gorilla/websocket"
	"github.com/vova616/screenshot"
)

var upgrader = websocket.Upgrader{}

func capture(c *websocket.Conn) error {
	screen, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, image.Image(screen), nil); err != nil {
		log.Println("unable to encode image.")
	}

	return c.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func captureHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("[read error]", err)
		return
	}
	if bytes.Equal(message, []byte("ok")) {
		log.Printf("received is not \"ok\", recv: %s", message)
		return
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("[read error]", err)
				return
			}
			log.Println("recv: ", message)
		}
	}()

	for {
		err = capture(c)
		if err != nil {
			log.Println("[write error]", err)
			break
		}

		time.Sleep(time.Second)
	}
}
