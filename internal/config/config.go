package config

import (
	"fmt"

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

    viper.AutomaticEnv()

    _ = viper.BindEnv("database.user", "DB_USER")
    _ = viper.BindEnv("database.password", "DB_PASSWORD")

	if err := viper.ReadInConfig(); err != nil {
        return Config{}, err
    }

	var cfg Config
	err := viper.Unmarshal(&cfg)
    fmt.Println("user:", cfg.Database.User, "password:", cfg.Database.Password)
	return cfg, err
}