package utils

import (
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()
