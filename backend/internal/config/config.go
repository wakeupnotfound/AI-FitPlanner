package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AI        AIConfig        `mapstructure:"ai"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Log       LogConfig       `mapstructure:"log"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	SecretKey string `mapstructure:"secret_key"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	PoolSize   int    `mapstructure:"pool_size"`
	MaxRetries int    `mapstructure:"max_retries"`
}

type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
}

type AIConfig struct {
	MaxConcurrentRequests int           `mapstructure:"max_concurrent_requests"`
	Timeout               time.Duration `mapstructure:"timeout"`
	RetryAttempts         int           `mapstructure:"retry_attempts"`
	RetryDelay            time.Duration `mapstructure:"retry_delay"`
}

type RateLimitConfig struct {
	APICallsPerMinute int64 `mapstructure:"api_calls_per_minute"`
	APICallsPerHour   int64 `mapstructure:"api_calls_per_hour"`
	APICallsPerDay    int64 `mapstructure:"api_calls_per_day"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

var GlobalConfig *Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 配置文件搜索路径
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/fitness-planner")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 绑定环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FITNESS")

	// 将配置解析到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	GlobalConfig = &config
	return nil
}

func setDefaults() {
	// 应用默认配置
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.name", "AI Fitness Planner")
	viper.SetDefault("app.version", "1.0.0")

	// 数据库默认配置
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.max_open_conns", 25)
	viper.SetDefault("database.mysql.max_idle_conns", 5)
	viper.SetDefault("database.mysql.conn_max_lifetime", "300s")

	viper.SetDefault("database.redis.port", 6379)
	viper.SetDefault("database.redis.db", 0)
	viper.SetDefault("database.redis.pool_size", 10)
	viper.SetDefault("database.redis.max_retries", 3)

	// JWT默认配置
	viper.SetDefault("jwt.access_token_expire", "3600s")
	viper.SetDefault("jwt.refresh_token_expire", "604800s")

	// AI默认配置
	viper.SetDefault("ai.max_concurrent_requests", 10)
	viper.SetDefault("ai.timeout", "60s")
	viper.SetDefault("ai.retry_attempts", 3)
	viper.SetDefault("ai.retry_delay", "5s")

	// 限流默认配置
	viper.SetDefault("rate_limit.api_calls_per_minute", 60)
	viper.SetDefault("rate_limit.api_calls_per_hour", 1000)
	viper.SetDefault("rate_limit.api_calls_per_day", 10000)

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("log.max_size", 500)
	viper.SetDefault("log.max_backups", 10)
	viper.SetDefault("log.max_age", 30)
}

func GetDSN() string {
	mysql := GlobalConfig.Database.MySQL
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.DBName)
}

func GetRedisAddr() string {
	redis := GlobalConfig.Database.Redis
	return fmt.Sprintf("%s:%d", redis.Host, redis.Port)
}
