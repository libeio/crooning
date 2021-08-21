
package xcache

import (
    "encoding/json"
    "log"
    "fmt"

    "github.com/go-redis/redis"

    "chat.com/xmodel"
)

const (
    _user_online_prefix    = "acc:user:online:"    // user online status
    _user_online_cache_time = 24 * 60 * 60
)

func _getUserOnlineKey(uk string) (string) {
    return fmt.Sprintf("%s%s", _user_online_prefix, uk)
}

func GetUserOnlineInfo(uk string) (*xmodel.UserOnline, error) {
    rc := getRedisCli()
    key := _getUserOnlineKey(uk)

    data, err := rc.Get(key).Bytes()
    if err != nil {
        if err == redis.Nil {
            log.Printf("-FATA- client(%s) redis.Do: %v", uk, err)
            return nil, err
        }
        log.Printf("-FATA- client(%s) redis.Do: %v", uk, err)
        return nil, err
    }

    uo := &xmodel.UserOnline{}
    err = json.Unmarshal(data, uo)
    if err != nil {
        log.Printf("-FATA- client(%s): %v", uk, err)
        return nil, err
    }

    log.Printf("-INFO- client(%s) online info: loginTime=%d,accIp=%s,isLogoff=%t", uk, uo.LoginTime, uo.AccIp, uo.IsLogoff)

    return uo, nil
}

func SetUserOnlineInfo(uk string, uo *xmodel.UserOnline) (error) {
    rc := getRedisCli()
    key := _getUserOnlineKey(uk)

    b, err := json.Marshal(uo)
    if err != nil {
        log.Printf("-FATA- client(%s): %v", uk, err)
        return err
    }

    _, err = rc.Do("setEx", key, _user_online_cache_time, string(b)).Result()
    if err != nil {
        log.Printf("-FATA- client(%s) redis.Do: %v", uk, err)
        return err
    }

    log.Printf("-INFO- client(%s) online update: loginTime=%d,accIp=%s,isLogoff=%t", uk, uo.LoginTime, uo.AccIp, uo.IsLogoff)

    return nil
}