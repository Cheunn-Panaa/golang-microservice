package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		App      appConf     `envPrefix:"APP_"`
		Server   serverConf  `envPrefix:"SERVER_"`
		Database DbConf      `envPrefix:"DATABASE_"`
		Swagger  swaggerConf `envPrefix:"SWAGGER_"`
		Logger   loggerConf  `envPrefix:"LOGGER_"`
	}
	appConf struct {
		Name string `env:"NAME,required" envDefault:"service"`
	}
	serverConf struct {
		Port         int           `env:"PORT,required"`
		TimeoutRead  time.Duration `env:"TIMEOUT_READ,required" envDefault:"60s"`
		TimeoutWrite time.Duration `env:"TIMEOUT_WRITE,required" envDefault:"60s"`
		TimeoutIdle  time.Duration `env:"TIMEOUT_IDLE,required" envDefault:"60s"`
	}

	DbConf struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,required"`
		Username string `env:"USER,required"`
		Password string `env:"PASSWORD,required"`
		DbName   string `env:"NAME,required"`
	}

	loggerConf struct {
		Level string `env:"LEVEL,required"`
	}

	swaggerConf struct {
		enable bool   `env:"ENABLE,required"`
		title  string `env:"TITLE,required"`
	}
)

func AppConfig() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return cfg
}

func DbConfig() *DbConf {
	var cfg DbConf
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	//log.debug("%+v\n", cfg)
	return &cfg
}
