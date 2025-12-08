package inits

import (
	"blog-post-service/rpc/models"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MysqlDb *gorm.DB
var configFile = flag.String("f", "etc/post.yaml", "the config file")

type DbInfo struct {
	MySQL Mysql
}
type Mysql struct {
	DSN           string
	IsAutoMigrate bool
}

// InitDB 初始化数据库连接
func InitDB() {
	// 初始化数据库连接
	// 这里可以使用数据库驱动的连接函数，如 mysql.Open()
	// 连接字符串格式："username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var c DbInfo
	conf.MustLoad(*configFile, &c)
	fmt.Println("mysql 信息为：", c.MySQL.DSN)
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
	db, err := gorm.Open(mysql.Open(c.MySQL.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移所有模型
	if c.MySQL.IsAutoMigrate {
		if err := db.AutoMigrate(
			&models.Post{},
		); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	}
	MysqlDb = db

	// 数据库连接成功
}
