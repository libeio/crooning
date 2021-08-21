
package xterminus

import (
    "sync"
    "net/http"

    "github.com/gorilla/websocket"
)

const (
    _default_app_id = 101   // default platform
)

var (
    _manager = NewManager()

    _app_ids = []uint32{ _default_app_id, 102, 103, 104 } // all the platforms
    
    _local_server_ip    string
    _local_server_port	string
)

type HandlerFunc func(*Client, string, []byte) (uint32, string, interface{})

var (
    _handlers = make(map[string]HandlerFunc)
    _handlers_rwmutex sync.RWMutex
)

var _upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}