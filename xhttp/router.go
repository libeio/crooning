
package xhttp

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "chat.com/xcode"
)

func InitRouter(router *gin.Engine) {
    router.LoadHTMLGlob("views/*")

    ur := router.Group("/user")
    {
        ur.GET("/list", 			userListFunc)
        ur.GET("/online",			userOnlineFunc)
        ur.POST("/sendMessage",		userSendMessageFunc)
        ur.POST("/sendMessageAll",	userSendMessageAllFunc)
    }

    sr := router.Group("/system")
    {
        sr.GET("/status", systemStatusFunc)
    }

    hr := router.Group("/home")
    {
        hr.GET("/index", homeIndexFunc)
    }
}

type codeNoticeData struct {
    Code    uint32      `json:"code"`
    Notice  string      `json:"notice"`
    Data    interface{} `json:"data"`
}

func httpResponse(c *gin.Context, code uint32, data map[string]interface{}) {
    notice := xcode.GetCodeNotice(code)
    cnd := codeNoticeData{
        Code:   code,
        Notice: notice,
        Data:   data,
    }

    // 允许跨域，这是允许访问所有域
	c.Header("Access-Control-Allow-Origin",      "*")
    // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
	c.Header("Access-Control-Allow-Methods",     "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
	c.Header("Access-Control-Allow-Headers",     "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	// 跨域关键设置 让浏览器可以解析
    c.Header("Access-Control-Expose-Headers",    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	//  跨域请求是否需要带cookie信息 默认设置为true
    c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, cnd)
}