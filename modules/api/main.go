package main

import (
	"flag"
	"fmt"
	"log"
	"spark-on-k8s-admin/config"
	"spark-on-k8s-admin/server"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config config_path")
	flag.Parse()
	cfg := config.Load(*configPath)
	fmt.Print(cfg)
	log.Default().Printf("Config: %+v", *cfg)
	server.Init(cfg)
}
