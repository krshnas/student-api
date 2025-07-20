package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// env-default:"production"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addr string
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH ")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configureation file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path not set")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}
	return &cfg
}
