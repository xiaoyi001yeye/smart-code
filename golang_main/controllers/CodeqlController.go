package controllers

import (
	"net/http"
	"github.com/kataras/iris/v12"
	"github.com/smartcodeql/services"
)

type CodeqlController struct {
	CodeQLService *services.CodeQLContainerService
}

func NewCodeqlController(service *services.CodeQLContainerService) *CodeqlController {
	return &CodeqlController{
		CodeQLService: service,
	}
}



func (ctrl *CodeqlController) GetContainerStatus(ctx iris.Context) {
	status, err := ctrl.CodeQLService.GetContainerStatus()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Writef("Error retrieving repositories: %v", err)
		return
	}

	ctx.JSON(status)
}