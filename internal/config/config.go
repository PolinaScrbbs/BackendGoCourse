package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Env         string     `yaml:"env" env-default:"local"`
    StoragePath string     `yaml:"storage_path" env-required:"true"`
    HTTPServer  HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
    Address     string        `yaml:"address" env-default:"localhost:8000"`
    Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set. Please specify the path to the configuration file.")
	}
	

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}
		log.Fatalf("error accessing config file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}