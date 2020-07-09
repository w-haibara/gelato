package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func shell() *os.File {
	c := exec.Command("bash")
	f, err := pty.Start(c)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return f
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	f := shell()
	defer f.Close()
	//	go io.Copy(os.Stdout, f)

	go func() {
		b := make([]byte, 1024)
		for {
			n, err := f.Read(b)
			if err != nil {
				log.Println(err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, b[:n])
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = f.Write(p)
		if err != nil {
			log.Println(err)
			return
		}

	}

}
