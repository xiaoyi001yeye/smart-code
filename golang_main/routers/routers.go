package routers

import (
    "github.com/kataras/iris/v12"
    "github.com/smartcodeql/services"
)

// ConfigureBooksRouter 设置书籍相关的路由
func ConfigureBooksRouter(app *iris.Application) {
    booksAPI :=app.Party("/books")
    {
        booksAPI.Use(iris.Compression)
        booksAPI.Get("/", services.GetBooks)
        booksAPI.Post("/", services.CreateBook)
        
    }

     
}