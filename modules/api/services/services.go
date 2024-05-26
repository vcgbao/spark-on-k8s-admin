package services

import (
	"spark-on-k8s-admin/config"
	"spark-on-k8s-admin/services/sparkoperator"
)

type Services struct {
	SparkOperatorService *sparkoperator.SparkOperatorService
}

func Init(cfg *config.Config) *Services {
	return &Services{
		SparkOperatorService: sparkoperator.Init(cfg),
	}
}
