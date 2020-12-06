package redislock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type LockItem struct {
	OwnerUnlockOnly bool   `json:"owner_unlock_only"` //是否只有自己可以解锁
	ClientId        string `json:"client_id"`
	Extra           string `json:"extra"`
}

type RedisLock struct {
	Client      *redis.Client
	ClientId    string //确保自己加的锁只有自己能解
	LockContext context.Context
}

func CreateRedisLock() (*RedisLock, error) {
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisLock{Client: client, LockContext: context.Background()}, nil
}

// 0表示这个锁key在解锁前永不过期
func (rl RedisLock) TryAcquire(key string, selfUnlockOnly bool,
	extra string, extraJudge func(extra string) error,
	expireTime time.Duration) (bool, error) {
	var lockItem = LockItem{
		OwnerUnlockOnly: selfUnlockOnly,
		ClientId:        rl.ClientId,
		Extra:           extra,
	}
	//这里需要一个等待机制（相当于锁定部分redis），一个额外的用来锁这个key的分布式锁
	//用来防止临界区抢占（即一个Client在TryAcquire中执行的时候有另一个Client也想加key相同的这个锁）
	//在这里，就是获取LockItem到重新设置LockItem这一段时间
	var cmd = rl.Client.SetNX(rl.LockContext, key, lockItem, expireTime)
	acquired, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return acquired, nil
}

//加锁（原语），不可重入
func (rl RedisLock) TryAcquireO(key string, clientId string, expireTime time.Duration) (bool, error) {
	var cmd = rl.Client.SetNX(rl.LockContext, key, clientId, expireTime)
	acquired, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return acquired, nil
}

//解锁（原语），不可重入
func (rl RedisLock) TryReleaseO(key string, clientId string) (bool, error) {
	var keyArgs = []string{key}
	//lua脚本作为一个事务，是原子性的
	var cmd = rl.Client.Eval(rl.LockContext, UNLOCK_LUA, keyArgs, clientId)
	result, err := cmd.Result()
	if err != nil {
		return false, err
	}
	var resultAsInt64 = result.(int64)
	return resultAsInt64 == 1, nil
}
