package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Host   string `mapstructure:"host"`
		Port   int    `mapstructure:"port"`
		ApiKey string `mapstructure:"api_key"`
	}
	Database struct {
		UserName string `mapstructure:"username"`
		PassWord string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	}
	ElasticSearch struct {
		Url      string `mapstructure:"url"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
}

var Conf = new(Config)

func Init() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	log.Println("配置文件加载成功")
}
