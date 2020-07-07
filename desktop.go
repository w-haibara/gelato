package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vova616/screenshot"
)

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

func desktopHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
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
	for !bytes.Equal(message, []byte("start")) {
		log.Printf("received is not \"ok\", recv: %s", message)
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

		time.Sleep(time.Millisecond * 30)
	}
}
