package redislock

const (
	UNLOCK_LUA = "if redis.call('get',KEYS[1]) == ARGV[1] then\n    return redis.call('del',KEYS[1])\nelse\n    return 0\nend"
)
