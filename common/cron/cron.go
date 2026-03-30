package cron

import (
	"exam-api-sync-go/service"
	"log"

	"github.com/robfig/cron/v3"
)

// Init 初始化定时任务
func Init() {
	log.Println("开始初始化定时任务...")
	c := cron.New()

	// 每天0点执行库存同步
	log.Println("添加库存同步定时任务...")
	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("开始执行库存同步任务...")
		syncService := service.NewInventorySyncService()
		if err := syncService.SyncInventory(); err != nil {
			log.Printf("库存同步失败: %v", err)
		} else {
			log.Println("库存同步成功")
		}
	})

	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
		return
	}

	// 启动定时任务
	log.Println("启动定时任务...")
	c.Start()
	log.Println("定时任务已启动")

	// 立即执行一次同步，确保数据初始化
	//log.Println("准备执行初始化库存同步...")
	//go func() {
	//	log.Println("进入初始化库存同步 goroutine...")
	//	time.Sleep(5 * time.Second)
	//	log.Println("执行初始化库存同步...")
	//	syncService := service.NewInventorySyncService()
	//	if err := syncService.SyncInventory(); err != nil {
	//		log.Printf("初始化库存同步失败: %v", err)
	//	} else {
	//		log.Println("初始化库存同步成功")
	//	}
	//}()
}
