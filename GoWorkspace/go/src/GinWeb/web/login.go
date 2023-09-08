package web

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	onChan    chan []byte
	inChan    chan []byte
	closeChan chan []byte
	isClosed  bool
	mutex     sync.Mutex
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		onChan:    make(chan []byte, 1000),
		inChan:    make(chan []byte, 1000),
		closeChan: make(chan []byte, 1),
	}
	// 启动读协程
	go conn.readLoop()

	// 启动写协程
	go conn.writerLoop()
	return
}

// 读取消息

func (conn *Connection) RedMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 发送消息

func (conn *Connection) WriterMessage(data []byte) (err error) {
	select {
	case conn.onChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 关闭连接

func (conn *Connection) Close() {
	conn.wsConn.Close()

	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()

}

// 循环读取消息写入conn.inChan

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}

	}
ERR:
	conn.Close()
}

// 循环发送消息写入conn.onChan

func (conn *Connection) writerLoop() {
	var (
		data []byte
		err  error
	)
	for {

		select {
		case data = <-conn.onChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
