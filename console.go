package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type SSHInfo struct {
	user string
}

func (info *SSHInfo) getSSHInfo(c *websocket.Conn) error {
	var (
		mt  int
		err error
	)
	err = c.WriteMessage(websocket.TextMessage, []byte("user: "))
	if err != nil {
		return err
	}
	user := []byte("")
	for {
		tmp := make([]byte, 1)
		mt, tmp, err = c.ReadMessage()
		if mt != websocket.TextMessage && err != nil {
			return err
		}
		err = c.WriteMessage(websocket.TextMessage, tmp[:1])
		if err != nil {
			return err
		}
		if tmp[0] == 0x0d {
			break
		}
		user = append(user, tmp[0])
	}
	err = c.WriteMessage(websocket.TextMessage, []byte("\n"))
	if err != nil {
		return err
	}
	info.user = string(bytes.Trim(user, " "))
	return nil
}

func (info *SSHInfo) shell() (*os.File, error) {
	c := exec.Command("/usr/bin/ssh", info.user+"@127.0.0.1")
	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	info := SSHInfo{}
	if err := info.getSSHInfo(c); err != nil {
		log.Println(err)
		return
	}

	f, err := info.shell()
	if f == nil {
		log.Println(info.user, ": user not found")
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

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
