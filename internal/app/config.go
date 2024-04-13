package app

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HTTPAddress      string        `yaml:"HTTP_ADDRESS" env:"HTTP_ADDRESS"`
	HTTPReadTimeout  time.Duration `yaml:"HTTP_READ_TIMEOUT" env:"PG_HOST"`
	HTTPWriteTimeout time.Duration `yaml:"HTTP_WRITE_TIMEOUT" env:"PG_HOST"`
	HTTPIdleTimeout  time.Duration `yaml:"HTTP_IDLE_TIMEOUT" env:"PG_HOST"`

	PgHost     string `yaml:"PG_HOST" env:"PG_HOST"`
	PgPassword string `yaml:"PG_PASSWORD" env:"PG_PASSWORD"`
	PgPort     string `yaml:"PG_PORT" env:"PG_PORT"`
	PgDBName   string `yaml:"PG_DBNAME" env:"PG_DBNAME"`
	PgUser     string `yaml:"PG_USER" env:"PG_USER"`
	PgSSLMode  string `yaml:"PG_SSLMODE" env:"PG_SSLMODE"`

	GrpcPort string `yaml:"GRPC_PORT" env:"GRPC_PORT"`
}

func MustLoadConfig(config_path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(config_path, &cfg)
	if err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	logrus.Debugln(cfg)
	return cfg, nil
}
