package main

import "github.com/kataras/iris/v12"
import "github.com/smartcodeql/routers"

func main() {
	app := iris.New()

	routers.ConfigureBooksRouter(app)


    root := func(ctx iris.Context) {
		ctx.Writef("Welcome to the Iris web server!")
	}

    app.Get("/", root)

	app.Listen(":3000")
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func list(ctx iris.Context) {
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	ctx.JSON(books)
	// 提示: 在服务器优先级和客户端请求中进行响应协商，
	// 以此来代替 ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func create(ctx iris.Context) {
	var b Book
	err := ctx.ReadJSON(&b)
	// 提示: 使用 ctx.ReadBody(&b) 代替，来绑定所有类型的入参
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Book creation failure").DetailErr(err))
		// 提示: 如果仅有纯文本（plain text）错误响应，
        // 可使用 ctx.StopWithError(code, err) 
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}