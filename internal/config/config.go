package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type Config struct {
	Server   ServerConfig      `mapstructure:"server"`
	Services map[string]string `mapstructure:"services"`
	Auth     AuthConfig        `mapstructure:"auth"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath) // Directly set the config file path

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("auth.jwt_secret", "")

	// Handle config file reading with proper error checking
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found: %v. Using defaults and environment variables.", err)
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	viper.AutomaticEnv() // Ensure env vars are read after config file

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Ensure Services map is initialized
	if config.Services == nil {
		config.Services = make(map[string]string)
	}

	return &config, nil
}
