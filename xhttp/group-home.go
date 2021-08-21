
package xhttp

import (
    "strconv"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"

    "chat.com/xterminus"
)

func homeIndexFunc(c *gin.Context) {
    str := c.Query("appId")
    i64, _ := strconv.ParseInt(str, 10, 32)
    appId := uint32(i64)

    if ! xterminus.ExistAppId(appId) {
        appId = xterminus.GetDefaultAppId()
    }

    // gin: type H map[string]interface{}
    data := gin.H {
        "title":        "chat-room-page",
        "appId":        appId,
        "httpUrl":      viper.GetString("app.httpUrl"),
        "webSocketUrl": viper.GetString("app.webSocketUrl"),
    }

    log.Printf("-INFO- home index: %v", data)

    c.HTML(http.StatusOK, "index.tpl", data)
}