package config

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type LogLevel string

const LogLevelDebug LogLevel = "debug"

const LogLevelError LogLevel = "error"

const LogLevelWarn LogLevel = "warn"

const LogLevelInfo LogLevel = "info"

type CalendarConfig struct {
	Logger  LoggerConf          `yaml:"logger" validate:"required"`
	Storage CalendarStorageConf `yaml:"storage" validate:"required"`
	HTTP    HTTPConf            `yaml:"http" validate:"required"`
	GRPC    GRPCConf            `yaml:"grpc" validate:"required"`
}

type SchedulerConfig struct {
	Logger  LoggerConf            `yaml:"logger" validate:"required"`
	Storage SchedulerStorageConf  `yaml:"storage" validate:"required"`
	Rabbit  SchedulerRabbitmqConf `yaml:"rabbit" validate:"required"`
	Cleanup CleanupConf           `yaml:"cleanup" validate:"required"`
}

type SenderConfig struct {
	Logger LoggerConf         `yaml:"logger" validate:"required"`
	Rabbit SenderRabbitmqConf `yaml:"rabbit" validate:"required"`
}

type CleanupConf struct {
	Days uint `yaml:"days" validate:"required"`
}

type HTTPConf struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type GRPCConf struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type CalendarStorageConf struct {
	Kind       string                 `yaml:"kind" validate:"required,oneof=db memory"`
	Connection DatabaseConnectionConf `yaml:"connection" validate:"required_if=Kind db"`
}

type SchedulerStorageConf struct {
	Connection DatabaseConnectionConf `yaml:"connection" validate:"required"`
}

type SchedulerRabbitmqConf struct {
	Connection RabbitmqConnectionConf `yaml:"connection" validate:"required"`
	Timer      RabbitmqTimerConf      `yaml:"timer" validate:"required"`
}

type SenderRabbitmqConf struct {
	Connection RabbitmqConsumerConf `yaml:"consumer" validate:"required"`
}

type RabbitmqTimerConf struct {
	Wait uint `yaml:"wait" validate:"required"`
}

type RabbitmqConnectionConf struct {
	Host     string `yaml:"host" validate:"required"`
	Port     uint16 `yalm:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Vhost    string `yaml:"vhost" validate:"required"`
	Exchange string `yaml:"exchange" validate:"required"`
}

type RabbitmqConsumerConf struct {
	Host     string `yaml:"host" validate:"required"`
	Port     uint16 `yalm:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Vhost    string `yaml:"vhost" validate:"required"`
	Exchange string `yaml:"exchange" validate:"required"`
	Queue    string `yaml:"queue" validate:"required"`
}

type LoggerConf struct {
	Level LogLevel `yaml:"level" validate:"required,oneof=debug error info"`
}

type DatabaseConnectionConf struct {
	Host     string `yaml:"host" validate:"required"`
	Port     uint16 `yalm:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Database string `yaml:"database" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func NewCalendarConfig(configFile string) (*CalendarConfig, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := CalendarConfig{}
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

func NewSchedulerConfig(configFile string) (*SchedulerConfig, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := SchedulerConfig{}
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

func NewSenderConfig(configFile string) (*SenderConfig, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := SenderConfig{}
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
