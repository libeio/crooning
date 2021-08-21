
package chat

import (
    "github.com/gin-gonic/gin"

    "chat.com/xhttp"
)

func InitRouter(router *gin.Engine) {
    xhttp.InitRouter(router)
}
