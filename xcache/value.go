
package xcache

// 服务器部分
const (
	// 保存集群中各服务器的哈希表名
	_servers_hash_key    		= "acc:hash:servers"
	// 哈希表名过期时间(秒计)
	_servers_hash_cache_time 	= 2 * 60 * 60
	// 服务器的超时时间(秒计)
	_servers_hash_timeout       = 3 * 60
)

// 客户端部分
const (
	// 用户连接超时
    _heartbeat_expiration_time = 6 * 60
)
