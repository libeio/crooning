
package xcode

const (
    BaseCode = 999
)

const (
    OK                 = 200            // Success
    NotLoggedIn        = BaseCode + 1   // 未登录
    ParameterIllegal   = BaseCode + 2   // 参数不合法
    UnauthorizedUserId = BaseCode + 3   // 非法的用户Id
    Unauthorized       = BaseCode + 4   // 未授权
    ServerError        = BaseCode + 5   // 系统错误
    NotData            = BaseCode + 6   // 没有数据
    ModelAddError      = BaseCode + 7   // 添加错误
    ModelDeleteError   = BaseCode + 8   // 删除错误
    ModelStoreError    = BaseCode + 9   // 存储错误
    OperationFailure   = BaseCode + 10  // 操作失败
    RoutingNotExist    = BaseCode + 11  // 路由不存在
)

// according to errcode project to the error information 
var _code_notice_map = map[uint32]string {
    OK:                 "Success",
    NotLoggedIn:        "未登录",
    ParameterIllegal:   "参数不合法",
    UnauthorizedUserId: "非法的用户Id",
    Unauthorized:       "未授权",
    NotData:            "没有数据",
    ServerError:        "系统错误",
    ModelAddError:      "添加错误",
    ModelDeleteError:   "删除错误",
    ModelStoreError:    "存储错误",
    OperationFailure:   "操作失败",
    RoutingNotExist:    "路由不存在",
}

func GetCodeNotice(code uint32) (string) {
    if notice, exists := _code_notice_map[code]; exists {
        return notice
    }
    return "未定义错误类型"
}