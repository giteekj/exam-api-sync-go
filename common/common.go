package common

import (
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestBody map[string]any

// GetRequestBody 获取Request请求Body体
func GetRequestBody(c *gin.Context) (body *RequestBody) {
	if c.Request.Method == "GET" {
		return body
	}
	// 将request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&body); err != nil {
		// 返回错误信息
		//c.JSON(http.StatusOK, common.Result(common.PARAMETER_ERROR, nil, err))
		//c.Abort()
		fmt.Println(err)
	}
	return body
}

func (c RequestBody) GetString(key string, context *gin.Context) string {
	v := c[key]
	if v == nil {
		return context.Query(key)
	}
	switch rv := v.(type) {
	case float64:
		return strconv.FormatFloat(rv, 'f', 0, 64)
	default:
		return v.(string)
	}
}

func (c RequestBody) GetInt(key string, context *gin.Context) int {
	v := c[key]
	if v == nil {
		i, _ := strconv.Atoi(context.Query(key))
		return i
	}
	return Int(v)
}

func Int(v any) int {
	if v == nil {
		return 0
	}
	switch rv := v.(type) {
	case int:
		return rv
	case int8:
		return int(rv)
	case int32:
		return int(rv)
	case uint32:
		return int(rv)
	case int64:

		return int(rv)
	case float64:
		return int(rv)
	case string:
		if rv == "" {
			return 0
		}
		i, _ := strconv.Atoi(rv)
		return i
	default:
		return 0
	}
}

func MD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func StrContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
