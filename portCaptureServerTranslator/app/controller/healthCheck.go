package controller

import (
	"github.com/gin-gonic/gin"
)

func (spc *sendPortController) HealthCheck(ctx *gin.Context) {

	err := spc.sendPortService.HealthCheck(ctx)
	if err != nil {
		ctx.AbortWithStatus(503)
		return
	}

	// ctx.JSON(200, getSuccessResponse())
}
