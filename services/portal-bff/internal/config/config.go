package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Database DBConfig     `mapstructure:"database"`
	Services map[string]string `mapstructure:"services"`
}

type ServerConfig struct {
	Name     string `mapstructure:"name"`
	HTTPPort int    `mapstructure:"http_port"`
}

type DBConfig struct {
	PortalDSN string `mapstructure:"portal_dsn"`
	IamDSN    string `mapstructure:"iam_dsn"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
