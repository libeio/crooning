
package xterminus

import (
    "time"
    "runtime/debug"
    "log"

    "chat.com/xcache"
)

type TimerFunc func(interface{}) (bool)

func taskTimer(delay, interval time.Duration, f TimerFunc, fParam interface{}, df TimerFunc, dfParam interface{}) {
    go func() {
        defer func() {
            if df != nil {
                df(dfParam)
            }
        } ()

        if f == nil {
            return
        }
    
        t := time.NewTimer(delay)
        defer t.Stop()
    
        for {
            select {
            case <- t.C:
                if f(fParam) == false {
                    return
                }
                t.Reset(interval)
            }
        }
    } ()
}

func TimerWebSocket() {
    // Note: interval time should not greater than _servers_hash_timeout
    taskTimer(2 * time.Second, 60 * time.Second, setServerInfoHandler, nil, delServerInfoHandler, nil)
}

func setServerInfoHandler(param interface{}) (bool) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    server := getLocalServer()

    now := time.Now()
    currentTime := uint64(now.Unix())

    // set the local server information to the cache server
    if err := xcache.SetServerInfo(server, currentTime); err != nil {
        log.Panicf("-FATA- set server cache: %s", err.Error())
    }

    log.Printf("-DBUG- local server: %s, currentTime: %d(%s)", server.String(), currentTime, now.Format(time.UnixDate))

    return true
}

func delServerInfoHandler(param interface{}) (bool) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    server := getLocalServer()
    if err := xcache.DelServerInfo(server); err != nil {
        log.Panicf("-FATA- del server cache error: %s", err.Error())
    }

    log.Printf("-INFO- del local server: %s", server.String())

    return false
}

func TimerCleanTimeoutConnections() {
    taskTimer(3 * time.Second, 30 * time.Second, cleanTimeoutConnectionsHandler, nil, nil, nil)
}

func cleanTimeoutConnectionsHandler(param interface{}) (bool) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf(string(debug.Stack()))
        }
    } ()

    CleanTimeoutConnections()

    log.Printf("-DBUG- clean timeout connections")

    return true
}