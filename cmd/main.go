package main

import (
	"exam-api-sync-go/common/cron"
	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/common/tools"
	"exam-api-sync-go/router"
	"fmt"
	"log"

	"exam-api-sync-go/common/redis"
	"exam-api-sync-go/common/setting"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	setting.Init()

	//初始化数据库连接
	orm.Init()

	// 初始化Redis连接
	redis.Init()

	// 初始化定时任务
	cron.Init()

	//直接执行一次库存同步
	//log.Println("直接执行库存同步任务...")
	//syncService := service.NewInventorySyncService()
	//if err := syncService.SyncInventory(); err != nil {
	//	log.Printf("库存同步失败: %v", err)
	//} else {
	//	log.Println("库存同步成功")
	//}

	fmt.Println("配置信息", *setting.GenFlag)

	if *setting.GenFlag != "" {
		tools.GenerateModel()

		fmt.Print("Generate Model Done.")
		return
	}

	// 设置Gin模式
	if setting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	router.RegisteAllRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", setting.Server.HttpPort)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
