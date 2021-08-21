
package xhttp

import (
    "strconv"
    "log"

    "github.com/gin-gonic/gin"

    "chat.com/xcode"
    "chat.com/xmodel"
    "chat.com/xcache"
    "chat.com/xterminus"
)

func userListFunc(c *gin.Context) {
    str := c.Query("appId")
    i64, _ := strconv.ParseInt(str, 10, 32)
    appId := uint32(i64)

    data := make(map[string]interface{})

    ul := xterminus.UserList(appId)
    data["userList"] = ul
    data["userCount"] = len(ul)

    log.Printf("-INFO- appId=%d, users count=%d", appId, len(ul))

    httpResponse(c, xcode.OK, data)
}

func userOnlineFunc(c *gin.Context) {
    userId := c.Query("userId")
    str := c.Query("appId")
    i64, _ := strconv.ParseInt(str, 10, 32)
    appId := uint32(i64)

    data := make(map[string]interface{})

    online := xterminus.CheckUserOnline(appId, userId)
    data["userId"] = userId
    data["online"] = online

    log.Printf("appId=%d, userId=%s, online=%t", appId, userId, online)

    httpResponse(c, xcode.OK, data)
}

func userSendMessageFunc(c *gin.Context) {
    str := c.PostForm("appId")
    i64, _ := strconv.ParseInt(str, 10, 32)
    appId := uint32(i64)
    userId := c.PostForm("userId")
    msgId := c.PostForm("msgId")
    message := c.PostForm("message")

    log.Printf("appId=%d, userId=%s, msgId=%s, message=%s", appId, userId, msgId, message)

    data := make(map[string]interface{})

    if xcache.SeqDuplicates(msgId) {
        log.Printf("repeated message(msgId=%d)!", msgId)
        httpResponse(c, xcode.OK, data)   // 重复发送消息
        return
    }

    bl, err := xterminus.SendMessageToUser(appId, userId, msgId, message)
    if err != nil {
        data["errstr"] = err.Error()
    }

    data["result"] = bl

    httpResponse(c, xcode.OK, data)
}

func userSendMessageAllFunc(c *gin.Context) {
    str := c.PostForm("appId")
    i64, _ := strconv.ParseInt(str, 10, 32)
    appId := uint32(i64)
    userId := c.PostForm("userId")
    msgId := c.PostForm("msgId")
    message := c.PostForm("message")

    data := make(map[string]interface{})

    if xcache.SeqDuplicates(msgId) {
        log.Printf("-WARN- repeated message(appId=%d,userId=%s,msgId=%s)", appId, userId, msgId)
        httpResponse(c, xcode.OK, data)
        return
    }

    log.Printf("-INFO- send to all: appId=%d,userId=%s,msgId=%s,message=%s", appId, userId, msgId, message)

    data["result"] = true
    err := xterminus.SendMessageToAll(appId, userId, msgId, xmodel.MessageCmdMsg, message)
    if err != nil {
        data["errstr"] = err.Error()
        data["result"] = false
    }

    httpResponse(c, xcode.OK, data)
}