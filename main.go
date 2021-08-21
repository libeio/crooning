
package main

import (
    "net/http"
    "log"

    "github.com/spf13/viper"
    "github.com/gin-gonic/gin"

    "chat.com/chat"
)

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

    viper.SetConfigName("conf/chat")
    viper.AddConfigPath(".")
    
    err := viper.ReadInConfig()
    if err != nil {
        panic("read conf error:" + err.Error())
    }
}

func main() {
    chat.InitRedis()

    router := gin.Default()

    chat.InitRouter(router)
    chat.InitEventHandler()

    chat.TimerWebSocket()
    chat.TimerCleanTimeoutConnections()

    go chat.StartWebs()
    go chat.StartGrpc()

    httpPort := viper.GetString("http.port")
    http.ListenAndServe(":" + httpPort, router)
}
