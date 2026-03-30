package orm

import (
	"fmt"
	"log"

	"exam-api-sync-go/common/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbs = make(map[string]*gorm.DB)

// Init 初始化数据库连接
func Init() {
	// 初始化fxshop数据库连接
	initDB("fxshop", *setting.Fxshop)

	// 初始化fxshop_sync数据库连接
	initDB("fxshop_sync", *setting.FxshopSync)
}

// initDB 初始化单个数据库连接
func initDB(name string, config setting.DatabaseConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Dbname)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库 %s 失败: %v", name, err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)

	dbs[name] = db
}

// GetDB 根据数据库名称获取数据库连接
func GetDB(name string) *gorm.DB {
	return dbs[name]
}
