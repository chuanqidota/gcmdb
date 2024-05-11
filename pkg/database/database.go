package database

import (
	"fmt"
	"gcmdb/config"
	"gcmdb/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Database.UserName,
		config.Conf.Database.PassWord,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.Name,
	)
	logger.Info("数据库连接信息:", dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Error("连接数据库出错:", err.Error())
		return
	}
	logger.Info("连接数据库成功：", dsn)
}
