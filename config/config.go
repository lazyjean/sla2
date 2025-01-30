package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	JWT      struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port    string
	Mode    string // 运行模式：debug, release
	Version string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	Debug           bool          // 是否开启调试模式
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 空闲连接最大生命周期
}

type RedisConfig struct {
	Host            string
	Port            string
	Password        string
	DB              int
	MaxRetries      int           // 最大重试次数
	MinRetryBackoff time.Duration // 最小重试间隔
	MaxRetryBackoff time.Duration // 最大重试间隔
	PoolSize        int           // 连接池大小
	MinIdleConns    int           // 最小空闲连接数
	MaxConnAge      time.Duration // 连接最大生命周期
}

type LogConfig struct {
	Level string // debug, info, warn, error
}

func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    "9000",
			Mode:    "debug",
			Version: "v1.0.0",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            "5432",
			User:            "sla",
			Password:        "sla1234",
			DBName:          "sla2",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 10 * time.Minute,
		},
		Redis: RedisConfig{
			Host:            "localhost",
			Port:            "6379",
			Password:        "",
			DB:              0,
			MaxRetries:      3,
			MinRetryBackoff: time.Millisecond * 100,
			MaxRetryBackoff: time.Second * 2,
			PoolSize:        100,
			MinIdleConns:    10,
			MaxConnAge:      30 * time.Minute,
		},
		Log: LogConfig{
			Level: "info", // 默认使用 info 级别
		},
		JWT: struct {
			SecretKey string `mapstructure:"secret_key"`
		}{
			SecretKey: "dj2m#9K$pL7&vX4@nR5*wQ8^hF3!tY6", // 32字节的随机字符串
		},
	}
}
