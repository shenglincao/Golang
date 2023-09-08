package service

import (
	"GinWeb/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func WebChat(c *gin.Context) {
	Chat(c.Writer, c.Request)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(w http.ResponseWriter, r *http.Request) {

	var (
		wsConn *websocket.Conn
		data   []byte
		err    error
		conn   *web.Connection
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	if conn, err = web.InitConnection(wsConn); err != nil {
		goto ERR
	}

	go func() {
		for {
			if err = conn.WriterMessage([]byte("hello")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.RedMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriterMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
