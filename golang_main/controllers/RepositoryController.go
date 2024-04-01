// controllers/RepositoryController.go

package controllers

import (
	// "encoding/json"
	"net/http"

	"github.com/kataras/iris/v12"
	// "github.com/smartcodeql/models"
	"github.com/smartcodeql/repositories"
)



// RepositoryController handles HTTP requests related to code repositories.
type RepositoryController struct {
	RepositoryDAO repositories.RepositoryDAO
}

// NewRepositoryController creates a new RepositoryController with the provided service.
func NewRepositoryController(dao repositories.RepositoryDAO) *RepositoryController {
	return &RepositoryController{
		RepositoryDAO: dao,
	}
}

// List handles the HTTP GET request for listing all repositories.
func (ctrl *RepositoryController) List(ctx iris.Context) {
	repositories, err := ctrl.RepositoryDAO.ListRepositories()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.Writef("Error retrieving repositories: %v", err)
		return
	}

	ctx.JSON(repositories)
}

// Get handles the HTTP GET request for getting a specific repository by ID.
func (ctrl *RepositoryController) Get(ctx iris.Context) {
	// repositoryID := ctx.Params().Get("id")
	// repository, err := ctrl.RepositoryDAO.GetRepositoryByID(repositoryID)
	// if err != nil {
	// 	// if errors.Is(err, repositories.ErrRepositoryNotFound) {
	// 	// 	ctx.StatusCode(http.StatusNotFound)
	// 	// 	ctx.Writef("Repository not found: %v", err)
	// 	// } else {
	// 	ctx.StatusCode(http.StatusInternalServerError)
	// 	ctx.Writef("Error retrieving repository: %v", err)
	// 	// }
	// 	return
	// }

	// ctx.JSON(repository)
}
