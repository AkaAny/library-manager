package rest

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Cache  bool        `json:"cache"`
}

func (br BaseResponse) Output(c *gin.Context) {
	c.JSON(br.Status, br)
}

type BaseError struct {
	Reason string      `json:"reason"`
	Extra  interface{} `json:"extra"`
}
