package rest

import (
	"github.com/NekoWheel/NekoCAS_Go_SDK"
	"github.com/gin-gonic/gin"
	"library-manager/logger"
	"net/http"
)

func nekoCasMiddleware() func(c *gin.Context) {
	cas := NekoCAS.New("127.0.0.1", "vNOZpKdqnUYcztBjUhvvPLpeYCIIBVev")
	return func(c *gin.Context) {
		var ticket = c.Query("ticket")
		logger.Info.Printf("cas ticket:%s", ticket)
		user, err := cas.Validate(ticket)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": err.Error(),
			}) //丢弃请求
			return
		}
		logger.Info.Printf("cas user:%v", user)
		c.Next()
	}
}
