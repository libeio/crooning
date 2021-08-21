
package xmodel

import (
    "log"
    "time"
)

const (
    heartbeatTimeout = 3 * 60
)

type UserOnline struct {
    AccIp         string `json:"accIp"`         // acc Ip
    AccPort       string `json:"accPort"`       // acc 端口
    AppId         uint32 `json:"appId"`         // appId
    UserId        string `json:"userId"`        // 用户Id
    ClientIp      string `json:"clientIp"`      // 客户端Ip
    ClientPort    string `json:"clientPort"`    // 客户端端口
    LoginTime     uint64 `json:"loginTime"`     // 用户上次登录时间
    HeartbeatTime uint64 `json:"heartbeatTime"` // 用户上次心跳时间
    LogoutTime    uint64 `json:"logoutTime"`    // 用户退出登录的时间
    Qua           string `json:"qua"`            // qua
    DeviceInfo    string `json:"deviceInfo"`    // 设备信息
    IsLogoff      bool   `json:"isLogoff"`      // 是否下线
}

func UserLogin(accIp, accPort string, appId uint32, userId string, addr string, loginTime uint64) (*UserOnline) {
    return &UserOnline{
        AccIp:		    accIp,
        AccPort:	    accPort,
        AppId:          appId,
        UserId:         userId,
        ClientIp:       addr,
        LoginTime:      loginTime,
        HeartbeatTime:  loginTime,
        IsLogoff:       false,
    }
}

func (this *UserOnline) Heartbeat(currentTime uint64) {
    this.HeartbeatTime = currentTime
    this.IsLogoff = false
}

func (this *UserOnline) Logout() {
    this.LogoutTime = uint64(time.Now().Unix())
    this.IsLogoff = true
}

func (this *UserOnline) IsOnline() bool {
    if this.IsLogoff {
        return false
    }

    currentTime := uint64(time.Now().Unix())
    if this.HeartbeatTime < (currentTime - heartbeatTimeout) {
        log.Printf("-ERRO- user(%d_%s) heartbeat(%d) timeout", this.AppId, this.UserId, this.HeartbeatTime)
        return false
    }

    return true
}

func (this *UserOnline) UserIsLocal(localIp, localPort string) bool {
    if this.AccIp == localIp && this.AccPort == localPort {
        return true
    }

    return false
}