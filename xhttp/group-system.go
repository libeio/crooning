
package xhttp

import (
    "runtime"
    "log"

    "github.com/gin-gonic/gin"

    "chat.com/xcode"
    "chat.com/xterminus"
)

func systemStatusFunc(c *gin.Context) {
    isDebug := c.Query("isDebug")

    data := make(map[string]interface{})

    numGoroutine := runtime.NumGoroutine()
    numCPU := runtime.NumCPU()

    data["numGoroutine"] = numGoroutine
    data["numCPU"] = numCPU

    data["managerInfo"] = xterminus.GetManagerInfo(isDebug)

    log.Printf("-INFO- system status(debug=%t): %v", isDebug, data)

    httpResponse(c, xcode.OK, data)
}