
package xterminus

import (
    "encoding/json"
    "runtime/debug"
    "log"

    "github.com/gorilla/websocket"

    "chat.com/xcode"
    "chat.com/xmodel"
)

const (
	// 用户连接超时时间
	_heartbeat_expiration_time = 6 * 60
)

/**
 * When client login to server, AppId and UserId must be needed for verifying.
 */
type login struct {
    AppId    uint32
    UserId   string
    Client  *Client
}

func (this *login) GetKey() string {
    return getUserKey(this.AppId, this.UserId)
}

// Client agent at server side
type Client struct {
    Addr 			string				// web client address
    Socket		   *websocket.Conn		// web client connection
    Send			chan []byte			// send message to peer websocket client
    AppId			uint32				// peer platform as logined, such as app/web/ios
    UserId			string				// will be assigned when user had logined
    FirstTime		uint64              // timestamp as first connect
    HeartbeatTime   uint64              // latest heartbeat time
    LoginTime       uint64              // timestamp when login
}

func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (*Client) {
    return &Client{
        Addr:           addr,
        Socket:         socket,
        Send:           make(chan []byte, 100),
        FirstTime:      firstTime,
        HeartbeatTime:  firstTime,
    }
}

func (this *Client) GetKey() (string) {
    return getUserKey(this.AppId, this.UserId)
}

func (this *Client) read() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    defer func() {
        close(this.Send)
        log.Printf("-NOTI- client(%s) close channel", this.Addr)
    } ()

    for {
        _, message, err := this.Socket.ReadMessage()
        if err != nil {
            log.Printf("-NOTI- client(%s) ReadMessage: %v", this.Addr, err)
            return
        }

        this.processData(message)
    }
}

func (this *Client) write() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    defer func() {
        _manager.Unregister <- this
        this.Socket.Close()
        log.Printf("-NOTI- client(%d_%s) close websocket connection", this.AppId, this.UserId)
    } ()

    for {
        select {
        case message, ok := <- this.Send:
            if ! ok {
                log.Printf("-NOTI- client(%d_%s) send failed and disconnect!", this.AppId, this.UserId)
                return
            }
            this.Socket.WriteMessage(websocket.TextMessage, message)
        }
    }
}

// 发送给 websocket 页面
func (this *Client) sendTo(msg []byte) {
    if this == nil {
        return
    }

    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    this.Send <- msg
}

func (this *Client) close() {
    close(this.Send)
}

func (this *Client) Login(appId uint32, userId string, loginTime uint64) {
    this.AppId = appId
    this.UserId = userId
    this.LoginTime = loginTime

    this.Heartbeat(loginTime)
}

func (this *Client) Heartbeat(currentTime uint64) {
    this.HeartbeatTime = currentTime
}

func (this *Client) IsHeartbeatTimeout(currentTime uint64) bool {
    if this.HeartbeatTime + _heartbeat_expiration_time <= currentTime {
        return true
    }
    return false
}

func (this *Client) IsLogin() (bool) {
    if this.UserId != "" {
        return true
    }
    return false
}

func (this *Client) processData(message []byte) {
    log.Printf("-DBUG- client(%s) received and process: %s", this.Addr, string(message))

    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    request := &xmodel.Request{}
    err := json.Unmarshal(message, request)
    if err != nil {
        log.Printf("-FATA- json.Unmarshal: %v", err)
        this.sendTo([]byte("message illegal"))
        return
    }

    requestData, err := json.Marshal(request.Data)
    if err != nil {
        log.Printf("-FATA- json.Marshal: %v", err)
        this.sendTo([]byte("request Data illegal"))
        return
    }

    seq := request.Seq
    cmd := request.Cmd

    log.Printf("-DBUG- client(%s) request: cmd=%s,seq=%s", this.Addr, cmd, seq)

    var (
        code uint32
        msg  string
        data interface{}
    )

    if handler, exists := getHandler(cmd); exists {
        code, msg, data = handler(this, seq, requestData)
    } else {
        code = xcode.RoutingNotExist
        log.Printf("-ERRO- router not exist with %s", cmd)
    }

    msg = xcode.GetCodeNotice(code)
    header := xmodel.NewResponseHeader(seq, cmd, code, msg, data)

    b, err := json.Marshal(header)
    if err != nil {
        log.Printf("-FATA- json.Marshal: %v", err)
        return
    }

    this.sendTo(b)

    log.Printf("-INFO- client(%d_%s) resp info: cmd=%s,code=%d", this.AppId, this.UserId, cmd, code)

    return
}