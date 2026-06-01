package config

import (
    "github.com/spf13/viper"
)

type ServerConfig struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     string `yaml:"port"`
    Mode     string `yaml:"mode"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
        return Config{}, err
    }

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return cfg, err
}