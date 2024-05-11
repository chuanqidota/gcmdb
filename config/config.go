package config

import (
	"gcmdb/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	Database struct {
		UserName string `json:"username"`
		PassWord string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
	}
}

var Conf = new(Config)

func Init() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("读取配置文件失败:", err.Error())
	}
	// 解析配置文件
	if err := viper.Unmarshal(Conf); err != nil {
		logger.Error("解析配置文件失败:", err.Error())
	}
	logger.Info("解析配置文件：", *Conf)
}
