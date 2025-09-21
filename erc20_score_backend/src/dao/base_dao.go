package dao

import (
	"time"

	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/config"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	// 连接到mysql数据库
	config := config.GetConfig()
	database := config.Database
	logger.Log.Printf("database is %v", database)
	db, err := gorm.Open(mysql.Open(database.User + ":" + database.Password +
		"@tcp(" + database.Host + ":" + database.Port + ")/" + database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		logger.Log.Printf("failed to connect database:%v", err)
	}

	// 设置连接池参数
	sqlDB, err := db.DB() // 获取*sql.DB对象进行更详细的配置
	if err != nil {
		logger.Log.Printf("failed to get sql db object:%s", err.Error())
	}
	sqlDB.SetMaxIdleConns(10)                  // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 设置最大打开的连接数
	sqlDB.SetConnMaxLifetime(60 * time.Second) // 设置连接可复用的最大时间（秒）

	return db

}
