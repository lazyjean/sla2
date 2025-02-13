package config

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//go:embed config-*.yaml
var configFS embed.FS

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Apple    AppleConfig    `mapstructure:"apple"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Debug           bool          `mapstructure:"debug"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	MaxRetries      int           `mapstructure:"max_retries"`
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxConnAge      time.Duration `mapstructure:"max_conn_age"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	TokenSecretKey   string `mapstructure:"token_secret_key"`
	RefreshSecretKey string `mapstructure:"refresh_secret_key"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

// AppleConfig 苹果登录配置
type AppleConfig struct {
	ClientID string `mapstructure:"client_id"`
}

var globalConfig *Config

func InitConfig() error {
	v := viper.New()
	v.SetConfigType("yaml")

	// 从环境变量获取环境配置，默认为 development
	env := os.Getenv("ACTIVE_PROFILE")
	var fileName string
	if env == "" {
		fileName = "config"
	} else {
		fileName = fmt.Sprintf("config-%s", env)
	}

	// try to read embedded config
	configFile := fmt.Sprintf("%s.yaml", fileName)
	fileData, err := configFS.ReadFile(configFile)
	if err == nil {
		if err := v.ReadConfig(bytes.NewReader(fileData)); err != nil {
			return fmt.Errorf("failed to read embedded config: %w", err)
		}
	} else {
		// 2. 如果嵌入文件不存在，尝试从本地文件系统读取
		v.SetConfigName(fileName)
		v.AddConfigPath(".") // 当前目录
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config from file system: %w", err)
		}
	}

	// 3. 支持从 Apollo 配置中心加载配置
	if v.GetBool("apollo.enabled") {
		apolloConfig := &config.AppConfig{
			AppID:          v.GetString("apollo.app_id"),
			Cluster:        v.GetString("apollo.cluster"),
			IP:             v.GetString("apollo.ip"),
			NamespaceName:  v.GetString("apollo.namespace"),
			IsBackupConfig: true,
			Secret:         v.GetString("apollo.secret"),
		}

		client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
			return apolloConfig, nil
		})
		if err != nil {
			return fmt.Errorf("failed to start Apollo client: %w", err)
		}

		// 将 Apollo 配置合并到 Viper 中
		apolloConfigMap := client.GetConfig(apolloConfig.NamespaceName).GetValue("content")
		if err := v.MergeConfig(bytes.NewBufferString(apolloConfigMap)); err != nil {
			return fmt.Errorf("failed to merge Apollo config: %w", err)
		}
	}

	// 自动加载环境变量
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.AllowEmptyEnv(true)

	// 设置配置监听
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := v.Unmarshal(&globalConfig); err != nil {
			fmt.Printf("Error reloading config: %v\n", err)
		}
	})

	if err := v.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func GetConfig() *Config {
	return globalConfig
}
