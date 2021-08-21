
package chat

import (
    "chat.com/xterminus"
)

func TimerWebSocket() {
    xterminus.TimerWebSocket()
}

func TimerCleanTimeoutConnections() {
    xterminus.TimerCleanTimeoutConnections()
}

func InitEventHandler() {
    xterminus.InitEventHandler()
}

func StartWebs() {
    xterminus.StartWebsServer()
}

func StartGrpc() {
    xterminus.StartGrpcServer()
}
