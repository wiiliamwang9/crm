package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config 配置结构体
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Static   StaticConfig   `yaml:"static"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
	Encoding string `yaml:"encoding"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
	WebPath string `yaml:"web_path"`
}

// StaticConfig 静态文件配置
type StaticConfig struct {
	EnableTargetRoute bool `yaml:"enable_target_route"`
	EnableConfigRoute bool `yaml:"enable_config_route"`
	EnableHealthCheck bool `yaml:"enable_health_check"`
}

// 全局变量
var (
	DB        *gorm.DB
	AppConfig *Config
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	AppConfig = &config
	return &config, nil
}

// ConnectDatabase 连接数据库
func ConnectDatabase() {
	if AppConfig == nil {
		log.Fatal("Config not loaded. Please call LoadConfig first.")
	}

	dbConfig := AppConfig.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s client_encoding=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
		dbConfig.SSLMode,
		dbConfig.TimeZone,
		dbConfig.Encoding,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")
	DB = database
}

// GetServerPort 获取服务器端口
func GetServerPort() string {
	if AppConfig == nil {
		return ":8081"
	}
	return fmt.Sprintf(":%d", AppConfig.Server.Port)
}

// SetGinMode 设置Gin模式
func SetGinMode() {
	if AppConfig != nil {
		switch AppConfig.Server.Mode {
		case "release":
			gin.SetMode(gin.ReleaseMode)
		case "test":
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode(gin.DebugMode)
		}
	}
}
