package rest

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/go-session/session"
	"library-manager/es"
	"library-manager/logger"
	"library-manager/utils"
	"library-manager/utils/redislock"
	"net/http"
)

func initDependencies() ([]interface{}, error) {
	var deps []interface{}
	deps = append(deps, "Library")
	//分布式锁
	redisLock, err := redislock.CreateRedisLock()
	if err != nil {
		return nil, err
	}
	deps = append(deps, redisLock)
	//ES
	les, err := es.CreateFromConfig()
	if err != nil {
		return nil, err
	}
	deps = append(deps, les)
	return deps, nil
}

func InitRestAPI() {
	var engine = gin.Default()
	engine.Use(ginsession.New())
	//对接NekoCAS
	var authController = AuthController{CASSecret: "http://127.0.0.1:8080/"}

	deps, err := initDependencies()
	if err != nil {
		panic(err)
	}
	engine.GET("/auth", utils.Warp(deps, authController.HandleAuth))
	engine.GET("/", utils.Warp(deps, func(c *gin.Context, strVal string) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"name": "ApiMain",
			//"decVal":decVal,
			"strVal": strVal,
		})
	}))

	var searchController = LibraryController{}
	var libGroup = engine.Group("/library")
	libGroup.Handle(http.MethodGet, "/search", utils.Warp(deps, searchController.HandleSearch))
	libGroup.Handle(http.MethodGet, "/esinfo", utils.Warp(deps,
		func(c *gin.Context, store session.Store, redisLock *redislock.RedisLock,
			les *es.LibraryES) {
			TryAcquireRedisLockByUserName(c, store, redisLock,
				func() {
					esInfo, err := les.GetInfo()
					if err != nil {
						c.JSON(http.StatusInternalServerError, BaseError{
							Reason: "fail to communicate with es",
							Extra:  err.Error(),
						})
						return
					}
					var response = BaseResponse{
						Status: http.StatusOK,
						Data:   esInfo.String(),
					}
					response.Output(c)
				})
		}))
	libGroup.Handle(http.MethodPost, "/add", utils.Warp(deps, searchController.HandleImport))

	err = engine.Run(":8080")
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
}
