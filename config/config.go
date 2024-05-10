package config

import (
	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Database        string `env:"DATABASE,notEmpty"`
	ElasticUrl      string `env:"ELASTIC_URL,notEmpty"`
	IndexName       string `env:"INDEX_NAME,notEmpty"`
	FilePath        string `env:"FILE_PATH,notEmpty"`
	DefaultPageSize int64  `env:"DEFAULT_PAGE_SIZE,notEmpty"`
	Log             LogConfig
}

type LogConfig struct {
	Output string `env:"OUTPUT" envDefault:"stdout"`
}

func Parse() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
