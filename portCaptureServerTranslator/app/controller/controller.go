package controller

import "github.com/gin-gonic/gin"

type SendPortsController interface {
	SendPorts(ctx *gin.Context)
}
