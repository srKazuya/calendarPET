package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Env string `yaml:"env" env-defaut:"dev"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct{
	Address string `yaml:"address" env-defaut:"0.0.0.0:8085"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}


func MustLoad() *Config  {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatal("CONFIG_PATH env is not set")
	}

	if _, err :=os.Stat(configPath); err != nil{
		log.Fatalf("err while oppening config file: %s", err)
	}

	var cfg Config
	
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("err while readig config, %s", err)
	}

	return &cfg
}