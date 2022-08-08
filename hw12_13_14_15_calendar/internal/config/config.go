package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type LogLevel string

const LogLevelDebug LogLevel = "debug"
const LogLevelError LogLevel = "error"
const LogLevelWarn LogLevel = "warn"
const LogLevelInfo LogLevel = "info"

type Config struct {
	Logger  LoggerConf  `yaml:"logger" validate:"required"`
	Storage StorageConf `yaml:"storage" validate:"required"`
}

type StorageConf struct {
	Kind string `yaml:"kind" validate:"required,oneof=db memory"`
	Dsn  string `yaml:"dsn" validate:"required_if=Kind db"`
}

type LoggerConf struct {
	Level LogLevel `yaml:"level" validate:"required,oneof=debug error info"`
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}
	v := validator.New()
	err = v.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to validate config: %w", err)
	}
	return &cfg, nil
}
