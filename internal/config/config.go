package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
	"time"
)

type Config struct {
	MaxWorkers    int           `env:"MAX_WORKERS"`
	MaxQueueSize  int           `env:"MAX_QUEUE_SIZE"`
	ServerAddress string        `env:"SERVER_ADDRESS"`
	ReadTimeout   time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout  time.Duration `env:"WRITE_TIMEOUT"`
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}
		if err := cleanenv.ReadEnv(config); err != nil {
			log.Fatalf("failed to parse configs: %v", err)
		}
	})

	return config, nil
}
