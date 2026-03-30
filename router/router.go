package router

import (
	"exam-api-sync-go/api"
	"exam-api-sync-go/middleware"

	"github.com/gin-gonic/gin"
)

// RegisteAllRoutes 注册所有路由
func RegisteAllRoutes(r *gin.Engine) {
	root := "/api/"
	router := r.Group(root)
	authorized := r.Group(root)

	// 授权中间件
	authorized.Use(middleware.JWTAuth())

	// 健康检查
	router.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 登录接口
	router.POST("login", api.Login)

	// 库存查询接口
	authorized.POST("inventory/query", api.InventoryQuery)

	// 其他路由可以在这里添加
}
