package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"library-manager/dbimport"
	"library-manager/es"
	"library-manager/logger"
	"library-manager/utils/redislock"
	"net/http"
)

type LibraryController struct {
}

const (
	MAX_CSV_SIZE = 1024 * 1024 * 2 //限制csv文件<=2MB
)

func getMatchCondition(c *gin.Context) es.H {
	var result = make(map[string]interface{})
	//var musts []interface{}
	var paramMap = c.Request.URL.Query()

	for k, v := range paramMap {
		result[k] = v[0]
	}

	return result
}

func (ls LibraryController) HandleSearch(c *gin.Context, les *es.LibraryES) {
	var result BaseResponse
	var title = c.Query("title")
	var author = c.Query("author")
	//c.Request.URL.Query()
	logger.Info.Printf("query book with title:%s author:%s", title, author)

	var matchMap = getMatchCondition(c)
	logger.Info.Printf("match map:%v", matchMap)

	var esBool = es.CreateByQueryMap(matchMap)
	resp, err := les.Search("marc", esBool)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Data = BaseError{
			Reason: "es search error",
			Extra:  err,
		}
		result.Output(c)
		return
	}
	var esResult gin.H
	err = json.NewDecoder(resp.Body).Decode(&esResult)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Data = BaseError{
			Reason: "es result unmarshall error",
			Extra:  err,
		}
		result.Output(c)
		return
	}
	result.Status = http.StatusOK
	result.Data = esResult
	result.Output(c)
}

func (ls LibraryController) HandleImport(c *gin.Context, store session.Store, redisLock *redislock.RedisLock,
	les *es.LibraryES) {
	TryAcquireRedisLockByUserName(c, store, redisLock, func() {
		csvHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, BaseError{
				Reason: "fail to get csv stream",
				Extra:  err.Error(),
			})
			return
		}
		logger.Info.Printf("csv file:%s %d", csvHeader.Filename, csvHeader.Size)
		if csvHeader.Size > MAX_CSV_SIZE { //csv流过大
			c.JSON(http.StatusBadRequest, BaseError{
				Reason: "csv data size cannot be bigger than 2MB",
			})
			return
		}
		handler, err := dbimport.CreateCSVToESHandler(les.GetClient())
		if err != nil {
			c.JSON(http.StatusInternalServerError, BaseError{
				Reason: "fail to create bulk session",
				Extra:  err.Error(),
			})
			return
		}
		err = handler.HandleCSVImport(csvHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, BaseError{
				Reason: "fail to add csv data to es",
				Extra:  err.Error(),
			})
			return
		}
		var response = BaseResponse{
			Status: http.StatusOK,
			Data:   "all lines imported",
		}
		response.Output(c)
	})

}
