package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"exam-api-sync-go/common"

	"github.com/gin-gonic/gin"
)

func Exception(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// buf := make([]byte, 1<<16)
			// runtime.Stack(buf, true)
			// fmt.Print("Exception:", r, string(buf))
			//打印错误堆栈信息
			log.Printf("panic-e: %v\n", r)
			debug.PrintStack()
			//封装通用json返回
			c.JSON(http.StatusOK, common.Result(common.SYSTEM_ERROR, nil, fmt.Sprintf("%v", r)))
		}
	}()
	//区分当前api和php版本（php 1 python 2 go 3）
	c.Header("v", "3")
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
