package rest

import (
	nekoCAS "github.com/NekoWheel/NekoCAS_Go_SDK"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"library-manager/logger"
	"library-manager/rest/model"
	"net/http"
)

type AuthController struct {
	CASSecret string
}

func (controller AuthController) HandleAuth(c *gin.Context, store session.Store) {
	var ticket = c.Query("ticket")
	logger.Info.Printf("cas ticket:%s", ticket)
	//domain参数是指CAS的endPoint
	cas := nekoCAS.New("http://127.0.0.1:8000", controller.CASSecret)
	user, err := cas.Validate(ticket)
	if err != nil {
		c.JSON(http.StatusUnauthorized, BaseError{
			Reason: "cas failed",
			Extra:  err.Error(),
		})
		return
	}
	store.Set("info", model.SessionInfo{
		Login:    true,
		UserName: user.Name,
	})
	err = store.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, BaseError{
			Reason: "fail to save session",
			Extra:  err.Error(),
		})
		return
	}
	var response = BaseResponse{
		Status: http.StatusOK,
		Data:   user,
	}
	response.Output(c)
}
