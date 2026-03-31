package setting

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	DatabaseConfigs = make(map[string]*DatabaseConfig)
	GenFlag         = flag.String("gen", "", "generate model file")
)

// Config 配置结构
type Config struct {
	Server     ServerConfig   `yaml:"server"`
	Fxshop     DatabaseConfig `yaml:"fxshop"`
	FxshopSync DatabaseConfig `yaml:"fxshop_sync"`
	Redis      RedisConfig    `yaml:"redis"`
	JWT        JWTConfig      `yaml:"jwt"`
	Sync       SyncConfig     `yaml:"sync"`
	Limit      LimitConfig    `yaml:"limit"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	RunMode         string `yaml:"run_mode"`
	HttpPort        int    `yaml:"http_port"`
	ApiKeyForRPC    string `yaml:"api_key_for_rpc"`
	ApiKeyForSMSRPC string `yaml:"api_key_for_sms_rpc"`
	GormLogLevel    int    `yaml:"gorm_log_level"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string   `yaml:"type"`
	Host            string   `yaml:"host"`
	Port            int      `yaml:"port"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	Dbname          string   `yaml:"dbname"`
	MaxIdleConn     int      `yaml:"max_idle_conn"`
	MaxOpenConn     int      `yaml:"max_open_conn"`
	ConnMaxLifetime int      `yaml:"conn_max_lifetime"`
	Replicas        []string `yaml:"replicas"`
	Sources         []string `yaml:"sources"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

// LimitConfig 限流配置
type LimitConfig struct {
	OssDownload [][]int `yaml:"oss_download"`
	OssData     [][]int `yaml:"oss_data"`
}

type SyncConfig struct {
	Token string `yaml:"token"`
}

// Global 全局配置变量
var (
	config     *Config
	Server     *ServerConfig
	Fxshop     *DatabaseConfig
	FxshopSync *DatabaseConfig
	Redis      *RedisConfig
	JWT        *JWTConfig
	Sync       *SyncConfig
)

// Init 初始化配置
func Init() {
	if !testing.Testing() {
		if !flag.Parsed() {
			flag.Parse()
		}
	}
	workDir, _ := os.Getwd()
	yamlFile, err := os.ReadFile(filepath.Join(workDir, "app.yaml"))
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	Server = &config.Server
	Fxshop = &config.Fxshop
	FxshopSync = &config.FxshopSync
	Redis = &config.Redis
	JWT = &config.JWT
	Sync = &config.Sync
	DatabaseConfigs[Fxshop.Dbname] = Fxshop
	DatabaseConfigs[FxshopSync.Dbname] = FxshopSync
}
