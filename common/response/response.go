package response

import (
	"exam-api-sync-go/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonResponse 返回正常数据
func JsonResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, common.Result(common.SUCCESS, data))
}

// ErrorResponse 返回错误数据
func ErrorResponse(c *gin.Context, code int, data any, msg ...any) {
	c.JSON(http.StatusOK, common.Result(code, data, msg...))
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}
