package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/smartcodeql/services"
)


type CodeqlController struct {
	codeqlContainerService services.CodeqlContainerService
}


func CodeqlController(service repositories.CodeqlContainerService) *CodeqlController {
	return &CodeqlController{
		codeqlContainerService: service,
	}
}


func (ctrl *CodeqlController) GetContainerStatus(ctx iris.Context) {
	status, err := ctrl.codeqlContainerService.GetContainerStatus()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Writef("Error retrieving repositories: %v", err)
		return
	}

	ctx.JSON(status)
}