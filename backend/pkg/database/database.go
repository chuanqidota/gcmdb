package database

import (
	"fmt"
	"gcmdb/config"
	"gcmdb/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Init() {
	host := envOrDefault("GCMDB_DATABASE_HOST", config.Conf.Database.Host)
	port := envOrDefault("GCMDB_DATABASE_PORT", fmt.Sprintf("%d", config.Conf.Database.Port))
	username := envOrDefault("GCMDB_DATABASE_USERNAME", config.Conf.Database.UserName)
	password := envOrDefault("GCMDB_DATABASE_PASSWORD", config.Conf.Database.PassWord)
	name := envOrDefault("GCMDB_DATABASE_NAME", config.Conf.Database.Name)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, name,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatalf("连接数据库出错: %v", err)
	}
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logger.Info("连接数据库成功")
}
