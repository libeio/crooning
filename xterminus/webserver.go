
package xterminus

import (
    "net/http"
    "time"
    "log"

    "github.com/spf13/viper"

    "chat.com/xmodel"
    "chat.com/xutil"
)

func getAppIds() []uint32 {
    return _app_ids
}

func getLocalServer() (*xmodel.Server) {
    return xmodel.NewServer(_local_server_ip, _local_server_port)
}

func isLocal(server *xmodel.Server) (bool) {
    if server.Ip == _local_server_ip && server.Port == _local_server_port {
        return true
    }
    return false
}

func ExistAppId(appId uint32) (bool) {
    for _, id := range _app_ids {
        if id == appId {
            return true
        }
    }
    return false
}

func GetDefaultAppId() (uint32) {
    return _default_app_id
}

func StartWebsServer() {
    _local_server_ip = xutil.GetServerIp()

    websocketPort := viper.GetString("websocket.port")
    _local_server_port = viper.GetString("rpc.port")

    http.HandleFunc("/acc", websocketUpgradeHandler)

    go _manager.Start()

    log.Printf("-INFO- websocket start on %s:%s", _local_server_ip, websocketPort)
    log.Printf("-INFO- internal communication port: %s", _local_server_port)

    http.ListenAndServe(":" + websocketPort, nil)
}

func websocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := _upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("-FATA- websocket Upgrade failed: %v", err)
        http.NotFound(w, r)
        return
    }
    // defer conn.Close()

    log.Printf("-INFO- websocket client(%s) upgrade!", conn.RemoteAddr().String())

    currentTime := uint64(time.Now().Unix())
    c := NewClient(conn.RemoteAddr().String(), conn, currentTime)

    go c.read()
    go c.write()

    _manager.Register <- c
}