package config

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppConfig   `mapstructure:",squash"`
	K8sConfig   `mapstructure:"k8s"`
	SparkConfig `mapstructure:"spark"`
}

func Load(path string) *Config {
	c, e := os.ReadFile(path)
	if e != nil {
		log.Fatalf("error while loading config file: %s", e)
	}
	var cfg = &Config{}
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(c))
	if err != nil {
		log.Fatal("Failed to read viper config", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Failed to unmarshal config", e)
	}

	fmt.Println(cfg.SparkConfig.Namespace)
	return cfg
}
