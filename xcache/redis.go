
package xcache

import (
    "log"

    "github.com/go-redis/redis"
    "github.com/spf13/viper"
)

var (
    _rc *redis.Client
)

func InitRedis() {
    _rc = redis.NewClient(&redis.Options{
        Addr:			viper.GetString("redis.addr"),
        Password:		viper.GetString("redis.password"),
        DB:				viper.GetInt("redis.DB"),
        PoolSize:		viper.GetInt("redis.poolSize"),
        MinIdleConns:   viper.GetInt("redis.minIdleConns"),
    })

    _, err := _rc.Ping().Result()
    if err != nil {
        log.Panicf("-FATA- connect redis failed")
    }
}

func getRedisCli() (*redis.Client) {
    return _rc
}