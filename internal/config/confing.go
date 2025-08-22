package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer config section
type HTTPServer struct {
	Addr string `yaml:"address" env:"HTTP_SERVER_ADDR" env-required:"true"`
}

// Config main structure
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

// MustLoad loads config or exits
func MustLoad() *Config {
	var configPath string

	if configPath = os.Getenv("CONFIG_PATH"); configPath == "" {
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()

		if configPath = *flags; configPath == "" {
			log.Fatal("Config Path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cfg
}
