package server

import (
	"fmt"
	"spark-on-k8s-admin/config"
	"spark-on-k8s-admin/controllers/proxy"
	"spark-on-k8s-admin/controllers/sparkoperator"
	"spark-on-k8s-admin/services"
	"spark-on-k8s-admin/utils"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, services *services.Services) *gin.Engine {
	sparkOperatorController := sparkoperator.Init(services.SparkOperatorService)

	router := gin.New()

	router.Use(
		static.Serve("/", static.LocalFile("/var/www/", false)),
	)

	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") && !strings.HasPrefix(c.Request.RequestURI, "/sparkui") {
			c.File("/var/www/index.html")
		}
	})

	sparkUiConfig := proxy.ApiConfig{
		SparkApplicationNamespace: cfg.SparkConfig.Namespace,
		SparkUIServiceUrl:         utils.DefaultIfBlank(cfg.SparkConfig.SparkUIServiceUrl, "http://{{$appName}}-ui-svc.{{$appNamespace}}.svc.cluster.local:4040"),
		ModifyRedirectUrl:         cfg.SparkConfig.ModifyRedirectUrl,
	}
	router.GET("/sparkui/*path",
		func(context *gin.Context) {
			proxy.ServeSparkUI(context, &sparkUiConfig, "/sparkui")
		})

	base := router.Group(fmt.Sprintf("%s/%s", cfg.BasePath, "/api"))

	base.GET("/sparkapplication", sparkOperatorController.GetSparkApplication)
	base.POST("/sparkapplication/create", sparkOperatorController.CreateSparkApplication)
	base.POST("/sparkapplication/update", sparkOperatorController.UpdateSparkApplication)
	base.DELETE("/sparkapplication/:name", sparkOperatorController.DeleteSparkApplication)

	base.GET("/scheduledsparkapplication", sparkOperatorController.GetScheduledSparkApplication)
	base.POST("/scheduledsparkapplication/create", sparkOperatorController.CreateScheduledSparkApplication)
	base.POST("/scheduledsparkapplication/update", sparkOperatorController.UpdateScheduledSparkApplication)
	base.DELETE("/scheduledsparkapplication/:name", sparkOperatorController.DeleteScheduledSparkApplication)

	return router
}
