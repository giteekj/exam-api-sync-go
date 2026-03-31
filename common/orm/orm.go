package orm

import (
	"fmt"
	"log"
	"os"
	"time"

	"exam-api-sync-go/common/setting"

	gormlog "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var dbs = make(map[string]*gorm.DB)

// Init 初始化数据库连接
func Init() {
	// 初始化fxshop数据库连接
	InitDB("fxshop", *setting.Fxshop, gormlog.LogLevel(setting.Server.GormLogLevel))

	// 初始化fxshop_sync数据库连接
	InitDB("fxshop_sync", *setting.FxshopSync, gormlog.LogLevel(setting.Server.GormLogLevel))
}

// InitDB 初始化单个数据库连接
func InitDB(name string, config setting.DatabaseConfig, logLevel logger.LogLevel) {

	// 初始化数据库日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n[Debug] ", log.Ldate|log.Ltime|log.Lshortfile), //
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Dbname)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("连接数据库 %s 失败: %v", name, err)
	}

	if len(config.Replicas) > 0 || len(config.Sources) > 0 {
		replicas := []gorm.Dialector{}
		for _, dns := range config.Replicas {
			replicas = append(replicas, mysql.Open(dns))
		}
		sources := []gorm.Dialector{}
		for _, dns := range config.Sources {
			sources = append(sources, mysql.Open(dns))
		}
		//读写分离
		//RegisterReplicas(setting.Fxshop, conn)
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Sources:  sources,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(int(config.MaxIdleConn)).
			SetMaxOpenConns(int(config.MaxOpenConn)).
			SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second))
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second) //设置连接空闲超时

	dbs[name] = db
}

// GetDB 根据数据库名称获取数据库连接
func GetDB(name string) *gorm.DB {
	return dbs[name]
}
