/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 全局配置变量，项目其他地方直接通过 config.AppConfig.Database.Host 调用
var AppConfig *Config

// Config 根配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Jwt      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	URL            string `mapstructure:"url"`
	FortuneQName   string `mapstructure:"fortuneQName"`
	EmbeddingQName string `mapstructure:"embeddingQName"`
	SearchQName    string `mapstructure:"searchQName"`
}

type LLMConfig struct {
	APIKey         string `mapstructure:"api_key"`
	BaseURL        string `mapstructure:"base_url"`
	Model          string `mapstructure:"model"`
	EmbeddingModel string `mapstructure:"embedding_model"`
}

type JWTConfig struct {
	Expire int    `mapstructure:"expire"`
	Secret string `mapstructure:"secret"`
	Prefix string `mapstructure:"prefix"`
}

// InitConfig 初始化配置文件
func InitConfig() {
	viper.SetConfigName("config")   // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")     // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config") // 查找配置文件的路径（相对于 main.go 启动路径）
	// 你也可以添加多个搜索路径，比如 viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// 将配置反序列化到全局变量 AppConfig 中
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("AppConfig无法decode到struct类型, 原因：%v", err)
	}

	// --- 关键：合并 prompt.yaml ---
	viper.SetConfigFile("config/prompt.yaml")
	if err := viper.MergeInConfig(); err != nil {
		zap.S().Errorf("加载 prompt.yaml 失败: %v", err)
	}

	fmt.Println("Config加载成功")
}

// GetPrompt 从配置中获取指定 key 的值，供其他模块调用
func GetPrompt(key string) string {
	return viper.GetString(key)
}
