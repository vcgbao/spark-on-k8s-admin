package server

import (
	"fmt"
	"spark-on-k8s-admin/config"
	"spark-on-k8s-admin/services"
)

func Init(cfg *config.Config) {
	services := services.Init(cfg)

	r := NewRouter(cfg, services)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
