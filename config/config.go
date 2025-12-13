package config

import (
	"os"
	"path/filepath"

	"shyft/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	App     App     `mapstructure:"app"`
	Auth    Auth    `mapstructure:"auth"`
	DB      DB      `mapstructure:"db"`
	Cache   Cache   `mapstructure:"cache"`
	Broker  Broker  `mapstructure:"broker"`
	Cookie  Cookie  `mapstructure:"cookie"`
	Session Session `mapstructure:"session"`
	Metric  Metric  `mapstructure:"metric"`
	Logger  Logger  `mapstructure:"logger"`
	Jaeger  Jaeger  `mapstructure:"jaeger"`
	Cdn     Cdn     `mapstructure:"cdn"`
}

type App struct {
	Mode    string `mapstructure:"mode"`
	Port    string `mapstructure:"port"`
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
}

type Auth struct {
	JwtPub string `mapstructure:"jwt_pub"`
}

type DB struct {
	Url string `mapstructure:"url"`
}

type Cache struct {
	Url string `mapstructure:"url"`
}

type Broker struct {
	Url           string `mapstructure:"url"`
	ConsumerGroup string `mapstructure:"consumer_group"`
	Topic         string `mapstructure:"topic"`
}

type Cookie struct {
	Name     string `mapstructure:"name"`
	MaxAge   int    `mapstructure:"max_age"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"http_only"`
}

type Session struct {
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Expire int    `mapstructure:"expire"`
}

type Metric struct {
	Url     string `mapstructure:"url"`
	Service string `mapstructure:"service"`
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}

type Jaeger struct {
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"service_name"`
	LogSpans    bool   `mapstructure:"log_spans"`
}

type Cdn struct {
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

var C Config

func ReadConfig(processCwdir string) {
	Config := &C
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(processCwdir, "config"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.CLogger.Errorf("Failed to read config file: %v", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		logger.CLogger.Errorf("Failed to unmarshal config file: %v", err)
		os.Exit(1)
	}

	logger.CLogger.Infof("Configuration loaded successfully from %s/config/.env.yaml", processCwdir)
}
