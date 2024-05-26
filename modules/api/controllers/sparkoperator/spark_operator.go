package sparkoperator

import (
	"net/http"
	"spark-on-k8s-admin/models"
	"spark-on-k8s-admin/services/sparkoperator"

	"github.com/gin-gonic/gin"
)

type SparkOperatorController struct {
	SparkOperatprService *sparkoperator.SparkOperatorService
}

func Init(SparkOperatprService *sparkoperator.SparkOperatorService) *SparkOperatorController {
	return &SparkOperatorController{
		SparkOperatprService: SparkOperatprService,
	}
}

func (c *SparkOperatorController) GetSparkApplication(ctx *gin.Context) {
	sparkApplicationResonse, err := c.SparkOperatprService.GetResource("sparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
	}
	ctx.JSON(
		http.StatusOK,
		models.Response[map[string]any]{
			StatusCode: models.SUCCESS,
			Data:       *sparkApplicationResonse,
		},
	)
}

func (c *SparkOperatorController) DeleteSparkApplication(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.SparkOperatprService.DeleteResource(name, "sparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)
}

func (c *SparkOperatorController) CreateSparkApplication(ctx *gin.Context) {
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	err := c.SparkOperatprService.CreateResource(body, "sparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)
}

func (c *SparkOperatorController) UpdateSparkApplication(ctx *gin.Context) {
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	err := c.SparkOperatprService.UpdateResource(body, "sparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)

}

func (c *SparkOperatorController) GetScheduledSparkApplication(ctx *gin.Context) {
	sparkApplicationResonse, err := c.SparkOperatprService.GetResource("scheduledsparkapplications")
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
	}
	ctx.JSON(
		http.StatusOK,
		models.Response[map[string]any]{
			StatusCode: models.SUCCESS,
			Data:       *sparkApplicationResonse,
		},
	)
}

func (c *SparkOperatorController) CreateScheduledSparkApplication(ctx *gin.Context) {
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	err := c.SparkOperatprService.CreateResource(body, "scheduledsparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)
}

func (c *SparkOperatorController) DeleteScheduledSparkApplication(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.SparkOperatprService.DeleteResource(name, "scheduledsparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)
}

func (c *SparkOperatorController) UpdateScheduledSparkApplication(ctx *gin.Context) {
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	err := c.SparkOperatprService.UpdateResource(body, "scheduledsparkapplications")

	if err != nil {
		ctx.JSON(
			http.StatusOK,
			models.Response[string]{
				StatusCode: models.ERROR,
				Message:    err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		models.Response[string]{
			StatusCode: models.SUCCESS,
		},
	)

}
