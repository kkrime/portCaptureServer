package controller

import (
	"context"
	"portCaptureServerTranslator/app/service"

	"github.com/gin-gonic/gin"
)

type sendPortController struct {
	sendPortService service.SendPortService
}

func NewSendPortController(sendPortService service.SendPortService) SendPortsController {
	return &sendPortController{
		sendPortService: sendPortService,
	}
}

func (spc *sendPortController) SendPorts(ctx *gin.Context) {

	portsData, err := ctx.GetRawData()
	if err != nil {
		abortAndError(ctx, err)
		return
	}

	err = spc.sendPortService.SendPort(context.Background(), &portsData)
	if err != nil {
		abortAndError(ctx, err)
		return
	}

	ctx.JSON(200, getSuccessResponse())
}
