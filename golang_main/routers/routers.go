package routers

import (
    "github.com/kataras/iris/v12"
    "github.com/smartcodeql/controllers"
    "github.com/smartcodeql/services"
)

// // ConfigureBooksRouter 设置书籍相关的路由
// func ConfigureBooksRouter(app *iris.Application) {
//     booksAPI :=app.Party("/books")
//     {
//         booksAPI.Use(iris.Compression)
//         booksAPI.Get("/", services.GetBooks)
//         booksAPI.Post("/", services.CreateBook)
        
//     }

     
// }


// func ConfigureRepoRouter(app *iris.Application) {
//     reposAPI :=app.Party("/repos")
//     {
//         reposAPI.Use(iris.Compression)
//         reposAPI.Get("/", controllers.List)
//         reposAPI.Get("/get", controllers.Get)
        
//     }

     
// }

func ConfigureCodeqlRouter(app *iris.Application) {
    service := services.NewCodeQLContainerService()
    codeQLController := controllers.NewCodeqlController(service)
    codeqlAPI :=app.Party("/codeql")
    {
        codeqlAPI.Use(iris.Compression)
        codeqlAPI.Get("/status", codeQLController.GetContainerStatus)
    }

     
}