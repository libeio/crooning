
package xterminus

import (
    "errors"
    "time"
    "log"

    "github.com/go-redis/redis"

    "chat.com/xcache"
    "chat.com/xmodel"
)

func UserList(appId uint32) ([]string) {
    ul := make([]string, 0)
    currentTime := uint64(time.Now().Unix())
    
    servers, err := xcache.GetAllServer(currentTime)
    if err != nil {
        log.Printf("GetAllServer: %v", err)
        return ul
    }

    for _, server := range servers {
        var list []string
        if isLocal(server) {
            list = getUserList(appId)
        } else {
            list, _ = grpcGetUserList(server, appId)
        }
        ul = append(ul, list...)
    }

    return ul
}

func CheckUserOnline(appId uint32, userId string) (isOnline bool) {
    if appId == 0 {
        for _, appId := range getAppIds() {
            isOnline, _ = checkUserOnline(appId, userId)
            if isOnline == true {
                break
            }
        }
    } else {
        isOnline, _ = checkUserOnline(appId, userId)
    }

    return isOnline
}

func checkUserOnline(appId uint32, userId string) (bool, error) {
    uk := getUserKey(appId, userId)
    uo, err := xcache.GetUserOnlineInfo(uk)
    if err != nil {
        if err == redis.Nil {
            log.Printf("-ERRO- client(%s) redis.Nil: %v", uk, err)
            return false, nil
        }
        log.Printf("-ERRO- client(%s) get online failed: %v", uk, err)
        return false, err
    }

    return uo.IsOnline(), nil
}

func SendMessageToUser(appId uint32, userId string, msgId, message string) (bool, error) {
    data := xmodel.GetTextMsgData(userId, msgId, message)

    // 位于本地机器上
    c := GetUserClient(appId, userId)
    if c != nil {
        bl, err := SendMessageToLocal(appId, userId, data)
        if err != nil {
            log.Printf("SendMessageToLocal: %v", err)
            return bl, err
        }
        return bl, nil
    }

    // 位于其他机器上(获取那台机器的缓存信息)
    uk := getUserKey(appId, userId)
    info, err := xcache.GetUserOnlineInfo(uk)
    if err != nil {
        log.Printf("GetUserOnlineInfo: %v", err)
        return false, err
    }

    // 将数据发往其他机器
    server := xmodel.NewServer(info.AccIp, info.AccPort)
    msg, err := grpcSendMsg(server, msgId, appId, userId, xmodel.MessageCmdMsg, xmodel.MessageCmdMsg, message)
    if err != nil {
        log.Printf("grpc.SendMsg(server=%s): %v", server.String(), err)
        return false, err
    }

    log.Printf("send message to user success. msg=%v", msg)
    
    return true, nil
}

func SendMessageToLocal(appId uint32, userId string, data string) (bool, error) {
    c := GetUserClient(appId, userId)
    if c == nil {
        return false, errors.New("not online")
    }

    c.sendTo([]byte(data))

    return true, nil
}

func SendMessageToAll(appId uint32, userId string, msgId, cmd, message string) (error) {
    currentTime := uint64(time.Now().Unix())
    servers, err := xcache.GetAllServer(currentTime)
    if err != nil {
        return err
    }

    for _, server := range servers {
        if isLocal(server) {
            data := xmodel.GetMsgData(userId, msgId, cmd, message)
            log.Printf("-DBUG- client(%d_%s) send: %v", appId, userId, data)
            sendMessageToAppId(appId, userId, data)
        } else {
            grpcSendMessageToAll(server, msgId, appId, userId, cmd, message)
        }
    }

    return nil
}
