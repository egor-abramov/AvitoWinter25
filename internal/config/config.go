package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"dev"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DB         `yaml:"db"`
	JWT        JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address     string        `yaml:"addr" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type DB struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     int    `yaml:"port" default:"5432"`
	User     string `yaml:"user" default:"postgres"`
	Database string `yaml:"database" default:"postgres"`
	Password string `env:"db_password" env-required:"true"`
}

type JWT struct {
	Secret  string        `env:"JWT_SECRET" env-required:"true"`
	ExpTime time.Duration `yaml:"exp_time" env-default:"24h"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "configs/dev.yml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist: %s", configPath))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
