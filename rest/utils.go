package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"library-manager/rest/model"
	"library-manager/utils/redislock"
	"net/http"
	"time"
)

func TryAcquireRedisLockByUserName(c *gin.Context,
	store session.Store, redisLock *redislock.RedisLock,
	doWhenAcquired func()) {
	var sessionInfo = model.GetFromStore(store)
	if sessionInfo == nil {
		c.JSON(http.StatusUnauthorized, BaseError{
			Reason: "add requires a logon user",
		})
		return
	}
	var lockKey = "es_lock"
	acquired, err := redisLock.TryAcquireO(lockKey, sessionInfo.UserName,
		10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, BaseError{
			Reason: "fail to communicate with redis",
			Extra:  err.Error(),
		})
		return
	}
	if !acquired {
		c.JSON(http.StatusBadRequest, BaseError{
			Reason: "lock is still existed",
		})
		return
	}
	//进入临界区
	doWhenAcquired()
	//释放分布式锁
	released, err := redisLock.TryReleaseO(lockKey, sessionInfo.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, BaseError{
			Reason: "fail to communicate with redis",
			Extra:  err.Error(),
		})
		return
	}
	if !released {

	}
}
