
package xcache

import (
    "log"
    "fmt"
)

const (
    _submit_again_prefix = "acc:submit:again:"
)

func _getSubmitAgainKey(from, value string) (string) {
    return fmt.Sprintf("%s%s:%s", _submit_again_prefix, from, value)
}

func _isSubmitAgain(from string, ttl int, value string) (bool) {
    key := _getSubmitAgainKey(from, value)

    rc := getRedisCli()
    number, err := rc.Do("setNx", key, "1").Int()
    if err != nil {
        log.Printf("-FATA- redis.Do: %v", err)
        return true
    }

    if number != 1 {
        return true
    }

    rc.Do("Expire", key, ttl)

    log.Printf("-INFO- from=%s add key=%s,value=%s, ttl to %d", from, key, "1", ttl)

    return false
}

func SeqDuplicates(seq string) (bool) {
    return _isSubmitAgain("seq", 12 * 60 * 60, seq)
}
