package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"production" env-required:"true"`
	HttpServer HTTPServer `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"5" env-required:"true"`
}

func MustLoad() *Config {
	configPath := "./config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cnf Config

	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("Can not read config file: %s", err)
	}

	return &cnf

}
