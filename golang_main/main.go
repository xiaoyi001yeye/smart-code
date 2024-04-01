package main

import "github.com/kataras/iris/v12"
import "github.com/smartcodeql/routers"

func main() {
	app := iris.New()

	routers.ConfigureCodeqlRouter(app)


    root := func(ctx iris.Context) {
		ctx.Writef("Welcome to the Iris web server!")
	}

    app.Get("/", root)

	app.Listen(":3000")
}