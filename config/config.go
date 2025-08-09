package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   Logger
}

// Logger config
type Logger struct {
	Development       bool   `env:LOG_DEVELOPMENT`
	DisableCaller     bool   `env:LOG_DISABLE_CALLER default:"false"`
	DisableStacktrace bool   `env:LOG_DISABLE_STACKTRACE default:"false"`
	Encoding          string `env:LOG_ENCODING`
	Level             string `env:LOG_LEVEL`
}
type ServerConfig struct {
	Port         string        `env:"PORT"`
	Mode         string        `env:"MODE"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
}
