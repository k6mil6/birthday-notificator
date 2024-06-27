package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env            string   `yml:"env" env-default:"local"`
	DB             DBConfig `yml:"db" env-required:"true"`
	HTTPPort       int      `yml:"http_port" env-default:"8080"`
	MigrationsPath string   `yml:"migrations_path" env-default:"./migrations"`
	JWT            JWT      `yml:"jwt" env-required:"true"`
}

type DBConfig struct {
	PostgresDSN   string        `yml:"postgres_dsn" env-default:"postgres://postgres:postgres@localhost:5442/birthday_notificator_db?sslmode=disable"`
	RetriesNumber int           `yml:"retries_number" env-default:"3"`
	RetryCooldown time.Duration `yml:"retry_cooldown" env-default:"10s"`
}

type JWT struct {
	Secret   string        `yml:"secret" env-required:"true"`
	TokenTTL time.Duration `yml:"token_ttl" env-default:"1h"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = "./config/config.yml" //default
	}

	return res
}
