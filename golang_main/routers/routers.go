package routers

import (
    "github.com/kataras/iris/v12"
    "github.com/smartcodeql/services"
    "github.com/smartcodeql/controllers"
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


func ConfigureRepoRouter(app *iris.Application) {
    reposAPI :=app.Party("/repos")
    {
        reposAPI.Use(iris.Compression)
        reposAPI.Get("/", controllers.List)
        reposAPI.Get("/get", controllers.Get)
        
    }

     
}

func ConfigureCodeqlRouter(app *iris.Application) {
    reposAPI :=app.Party("/codeql")
    {
        reposAPI.Use(iris.Compression)
        reposAPI.Get("/status", controllers.GetContainerStatus)
    }

     
}