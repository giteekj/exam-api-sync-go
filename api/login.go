package api

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/model/request"
	"net/http"

	"exam-api-sync-go/common/response"
	"exam-api-sync-go/service"

	"github.com/gin-gonic/gin"
)

// Login 登录接口
func Login(c *gin.Context) {
	body := common.GetRequestBody(c)
	if body == nil {
		response.ErrorResponse(c, common.FATAL, nil)
		return
	}
	username := body.GetString("username", c)
	password := body.GetString("password", c)
	if username == "" || password == "" {
		c.JSON(http.StatusOK, common.Result(common.LACK_OF_PARAMETER, nil))
		return
	}
	// 调用登录服务
	loginService := service.NewLoginService()
	result, err := loginService.Login(request.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		response.ErrorResponse(c, common.FATAL, nil, err.Error())
		return
	}

	response.JsonResponse(c, result)
}
