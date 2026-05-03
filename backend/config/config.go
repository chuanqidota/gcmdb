package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		AutoRelation  bool   `mapstructure:"auto_relation"`
		SessionMaxAge int    `mapstructure:"session_max_age"`
	}
	Database struct {
		UserName string `mapstructure:"username"`
		PassWord string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	}
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
	} `mapstructure:"cors"`
	TokenCacheTTL int `mapstructure:"token_cache_ttl"` // 单位：分钟，默认 5
}

var Conf = new(Config)

func Init() {
	// 支持 GCMDB_CONFIG 环境变量覆盖配置文件路径
	configPath := os.Getenv("GCMDB_CONFIG")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}
	viper.SetConfigFile(configPath)
	viper.SetEnvPrefix("GCMDB")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	log.Println("配置文件加载成功")
}
