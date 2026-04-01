package main

import (
	"exam-api-sync-go/common/cron"
	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/common/tools"
	"exam-api-sync-go/middleware"
	"exam-api-sync-go/router"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"exam-api-sync-go/common/redis"
	"exam-api-sync-go/common/setting"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var start = time.Now()

func stepCost(step int) {
	tc := time.Since(start)
	fmt.Printf("执行%d耗时： %v\n", step, tc)
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			// buf := make([]byte, 1<<16)
			// runtime.Stack(buf, true)
			// fmt.Print("Exception:", r, string(buf))
			//打印错误堆栈信息
			log.Printf("panic0: %v\n", r)
			debug.PrintStack()
		}
	}()

	// 初始化配置
	setting.Init()
	stepCost(0)

	//初始化数据库连接
	orm.Init()
	stepCost(1)

	// 初始化Redis连接
	redis.Init()
	stepCost(2)

	// 初始化定时任务
	cron.Init()
	stepCost(3)

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
	stepCost(5)

	// 创建Gin实例
	engine := gin.New()
	//engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.RedirectFixedPath = true
	//添加全局异常处理，否则任何空指针/内存异常会导致蓝屏
	engine.Use(middleware.Exception)
	engine.Use(middleware.Cors)
	//engine.Use(middleware.RouterPathIgnoreCase)
	engine.Use(gzip.Gzip(gzip.BestSpeed))
	stepCost(6)

	// 注册路由
	router.RegisteAllRoutes(engine)
	stepCost(7)

	// 启动服务器
	addr := fmt.Sprintf(":%d", setting.Server.HttpPort)
	stepCost(8)
	log.Printf("服务器启动在 %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
