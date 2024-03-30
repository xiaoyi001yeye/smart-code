package services

import (
    "github.com/kataras/iris/v12"
)

// Book 结构体定义书籍信息
type Book struct {
    Title string `json:"title"`
}

// BooksList 用于存储书籍列表
var BooksList = []Book{
    {Title: "Mastering Concurrency in Go"},
    {Title: "Go Design Patterns"},
    {Title: "Black Hat Go"},
}

// GetBooks 处理 GET 请求，返回书籍列表
func GetBooks(ctx iris.Context) {
    ctx.JSON(BooksList)
}

// CreateBook 处理 POST 请求，创建新的书籍记录
// 这里仅为示例，实际应用中应该包含更完整的逻辑
func CreateBook(ctx iris.Context) {
    // 从请求中解析书籍信息
    // 这里省略了错误处理和实际的数据库操作
    book := Book{Title: "New Book Title"}
    BooksList = append(BooksList, book)
    ctx.JSON(iris.Map{"message": "Book created successfully"})
}