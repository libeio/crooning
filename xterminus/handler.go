
package xterminus

import (
    "time"
    "encoding/json"
    "log"

    "github.com/go-redis/redis"

    "chat.com/xcode"
    "chat.com/xmodel"
    "chat.com/xcache"
)

func InitEventHandler() {
    register("ping",      pingHandlerFunc)
    register("login",     loginHandlerFunc)
    register("heartbeat", heartbeatHandlerFunc)
}

func register(cmd string, h HandlerFunc) {
    _handlers_rwmutex.Lock()
    defer _handlers_rwmutex.Unlock()
    _handlers[cmd] = h
}

func getHandler(cmd string) (HandlerFunc, bool) {
    _handlers_rwmutex.RLock()
    defer _handlers_rwmutex.RUnlock()
    h, exists := _handlers[cmd]

    return h, exists
}

func pingHandlerFunc(c *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
    code = xcode.OK
    data = "pong"

    log.Printf("client(%d_%s) seq=%s message=%s", c.AppId, c.UserId, seq, string(message))

    return
}

func loginHandlerFunc(c *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
    code = xcode.OK
    currentTime := uint64(time.Now().Unix())

    request := &xmodel.Login{}
    if err := json.Unmarshal(message, request); err != nil {
        code = xcode.ParameterIllegal
        log.Printf("-FATA- client(%s): %v", c.Addr, err)
        return
    }

    log.Printf("-INFO- client(%s): seq=%s,ServiceToken=%s,login=%s", c.Addr, seq, request.ServiceToken, string(message))

    if request.UserId == "" || len(request.UserId) >= 20 {
        code = xcode.UnauthorizedUserId
        log.Printf("-ERRO- client(%s) illegal UserId: %s", c.Addr, request.UserId)
        return
    }

    if ! ExistAppId(request.AppId) {
        code = xcode.Unauthorized
        log.Printf("-ERRO- client(%s) AppId not exist: %d", c.Addr, request.AppId)
        return
    }

    if c.IsLogin() {
        code = xcode.OperationFailure
        log.Printf("-WARN- client(%s) already logined", c.Addr)
        return
    }

    c.Login(request.AppId, request.UserId, currentTime)

    // store in the cache
    userOnline := xmodel.UserLogin(_local_server_ip, _local_server_port, request.AppId, request.UserId, c.Addr, currentTime)
    err := xcache.SetUserOnlineInfo(c.GetKey(), userOnline)
    if err != nil {
        code = xcode.ServerError
        log.Printf("-ERRO- client(%s) set online: %v", c.Addr, err)
        return
    }

    login := &login{
        AppId:  request.AppId,
        UserId: request.UserId,
        Client: c,
    }

    // send the login information to manager
    _manager.Login <- login

    log.Printf("-INFO- client(%s,%d_%s) logined success", c.Addr, c.AppId, c.UserId)

    return
}

func heartbeatHandlerFunc(c *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
    code = xcode.OK
    currentTime := uint64(time.Now().Unix())
    
    request := &xmodel.HeartBeat{}
    if err := json.Unmarshal(message, request); err != nil {
        code = xcode.ParameterIllegal
        log.Printf("-FATA- client(%d_%s): %v", c.AppId, c.UserId, err)
        return
    }

    log.Printf("-DBUG- client(%d_%s) heartbeat: seq=%s,message=%s", c.AppId, c.UserId, seq, string(message))

    if ! c.IsLogin() {
        code = xcode.NotLoggedIn
        log.Printf("-ERRO- client(%d_%s) not logined", c.AppId, c.UserId)
        return
    }

    userOnline, err := xcache.GetUserOnlineInfo(c.GetKey())
    if err != nil {
        if err == redis.Nil {
            code = xcode.NotLoggedIn
            log.Printf("-ERRO- client(%d_%s) redis Nil", c.AppId, c.UserId)
            return
        } else {
            code = xcode.ServerError
            log.Printf("-ERRO- client(%d_%s) get online: %v", c.AppId, c.UserId, err)
            return
        }
    }

    c.Heartbeat(currentTime)
    userOnline.Heartbeat(currentTime)
    err = xcache.SetUserOnlineInfo(c.GetKey(), userOnline)
    if err != nil {
        code = xcode.ServerError
        log.Printf("-ERRO- client(%d_%s) set online: %v", c.AppId, c.UserId, err)
        return
    }

    return
}