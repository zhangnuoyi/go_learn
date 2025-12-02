package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"04blog/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	//服务器配置
	Server ServerConfig
	//数据库配置
	DB DBConfig
	//Redis配置
	Redis RedisConfig
}

type ServerConfig struct {
	Port string //服务器端口
	Host string //服务器主机
}

type DBConfig struct {
	Host          string //数据库主机
	Port          string //数据库端口
	User          string //数据库用户名
	Pass          string //数据库密码
	Name          string //数据库名
	IsAutoMigrate bool   //是否自动迁移数据库
}

type RedisConfig struct {
	Host     string //Redis主机
	Port     string //Redis端口
	Password string //Redis密码
	DB       int    //Redis数据库索引
}

// AppConfig 全局配置实例
var AppConfig Config

// RedisClient 全局Redis客户端实例
var RedisClient *redis.Client

func LoadConfig() error {
	//走环境变量加载配置信息  如果没有 则走默认值
	AppConfig = Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		DB: DBConfig{
			Host:          "172.18.112.82", //wls 虚拟化ip  数据库主机
			Port:          "3306",
			User:          "root",
			Pass:          "root",
			Name:          "moon_blog",
			IsAutoMigrate: false,
		},
		Redis: RedisConfig{
			Host:     "172.18.112.82",
			Port:     "6379",
			Password: "123456",
			DB:       0,
		},
	}
	return nil
}

func InitDB() (*gorm.DB, error) {
	// 初始化数据库连接
	// 这里可以使用数据库驱动的连接函数，如 mysql.Open()
	// 连接字符串格式："username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DB.User, AppConfig.DB.Pass,
		AppConfig.DB.Host, AppConfig.DB.Port,
		AppConfig.DB.Name)
	//数据库操作 日志设置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移所有模型
	if AppConfig.DB.IsAutoMigrate {
		if err := db.AutoMigrate(
			&models.User{},
			&models.Post{},
			&models.Comment{},
			&models.Like{},
		); err != nil {
			panic("failed to migrate database")
		}
	}

	// 数据库连接成功
	return db, nil
}

// InitRedis 初始化Redis连接
func InitRedis() (*redis.Client, error) {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppConfig.Redis.Host, AppConfig.Redis.Port),
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v\n", err)
		return nil, err
	}

	log.Println("Redis连接成功")
	return RedisClient, nil
}
