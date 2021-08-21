
package xcache

import (
    "strconv"
    "log"
    "fmt"

    "chat.com/xmodel"
)

func _getServerHashKey() (string) {
    return fmt.Sprintf("%s", _servers_hash_key)
}

func SetServerInfo(server *xmodel.Server, timestamp uint64) (error) {
    hashKey := _getServerHashKey()
    timestampString := fmt.Sprintf("%d", timestamp)

    rc := getRedisCli()

    number, err := rc.Do("hSet", hashKey, server.String(), timestampString).Int()
    if err != nil {
        return err
    }

    if number != 1 {
        ttl, err := rc.Do("ttl", hashKey).Int()
        if err == nil {
            log.Printf("-NOTI- field=%s already in hash=%s, ttl=%d", server.String(), hashKey, ttl)
            return nil
        } else {
            return err
        }
    }

    rc.Do("Expire", hashKey, _servers_hash_cache_time)

    log.Printf("-INFO- add field=%s,value=%s to hash=%s, reset ttl to %d", server.String(), timestampString, hashKey, _servers_hash_cache_time)

    return nil
}

func DelServerInfo(server *xmodel.Server) (error) {
    hashKey := _getServerHashKey()

    rc := getRedisCli()

    number, err := rc.Do("hDel", hashKey, server.String()).Int()
    if err != nil {
        return err
    }

    if number != 1 {
        return nil
    }

    rc.Do("Expire", hashKey, _servers_hash_cache_time)

    log.Printf("-INFO- del field=%s from hash=%s, reset ttl to %d", server.String(), hashKey, _servers_hash_cache_time)

    return nil
}

func GetAllServer(currentTime uint64) ([]*xmodel.Server, error) {
    servers := make([]*xmodel.Server, 0)
    hashKey := _getServerHashKey()

    rc := getRedisCli()

    serverMap, err := rc.HGetAll(hashKey).Result()
    if err != nil {
        return nil, err
    }

    for k, v := range serverMap {
        u64, err := strconv.ParseUint(v, 10, 64)
        if err != nil {
            return nil, err
        }

        // timeout _servers_hash_timeout should greater than taskTimer
        if u64 + _servers_hash_timeout <= currentTime {
            continue
        }

        server, err := xmodel.ToServer(k)
        if err != nil {
            return nil, err
        }

        servers = append(servers, server)
    }

    log.Printf("-INFO- all available server: %v", servers)

    return servers, nil
}