package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	DB             DBConfig      `yaml:"db" env-required:"true"`
	HTTPPort       int           `yaml:"http_port" env-default:"8080"`
	MigrationsPath string        `yaml:"migrations_path" env-default:"./migrations"`
	JWT            JWT           `yaml:"jwt" env-required:"true"`
	Email          Email         `yaml:"email" env-required:"true"`
	ScanInterval   time.Duration `yaml:"scan_interval" env-default:"1m"`
}

type DBConfig struct {
	PostgresDSN   string        `yaml:"postgres_dsn" env-default:"postgres://postgres:postgres@localhost:5442/birthday_notificator_db?sslmode=disable"` //nolint:lll
	RetriesNumber int           `yaml:"retries_number" env-default:"3"`
	RetryCooldown time.Duration `yaml:"retry_cooldown" env-default:"10s"`
}

type JWT struct {
	Secret   string        `yaml:"secret" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type Email struct {
	SenderAddress  string `yaml:"sender_address" env-required:"true"`
	SenderPassword string `yaml:"sender_password" env-required:"true"`
	SMTPAddress    string `yaml:"smtp_address" env-required:"true"`
	SMTPPort       int    `yaml:"smtp_port" env-default:"587"`
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
		res = "./config/config.yaml" //default
	}

	return res
}
