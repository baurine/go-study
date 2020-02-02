# Go Misc

- Go Modules
- Web 框架
- Go 中的字符串

## Go Modules

- [Go Modules 详解使用](https://learnku.com/articles/27401)
- [用 go-module 作为包管理器搭建 go 的 web 服务器](https://www.hulunhao.com/go/go-web-backend-starter/)

go 官方依赖管理。`go mod init module_name`，生成 go.mod 文件，执行 `go build` 会生成 go.sum，go.mod 相当于 npm 中的 package.json，go.sum 相当于 package-lock.json。

示例：

1. 在 GOPATH 以外的路径下，`mkdir backend && cd backend`
1. 执行 `go mod init backend` 初始化，生成 go.mod，其内容如下

   ```
   module backend

   go 1.12
   ```

1. 创建 main.go，import gin 包

   ```go
   package main

   import "github.com/gin-gonic/gin"

   func main() {
       r := gin.Default()
       ...
   }
   ```

1. 执行 `go build`，自动下载依赖包，go.mod 被自动更新并生成 go.sum

   ```
   // go.mod
   module backend

   go 1.12

   require(
       github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
       github.com/gin-gonic/gin v1.3.0
       ...
   )
   ```

1. 查看依赖：`go list -m all`

1. 清除过期依赖：`go mod tidy`

1. 更多命令略

## Web 框架

资源：

- [7 天用 Go 从零实现 Web 框架 Gee 教程](https://github.com/geektutu/7days-golang)
- [Go-Mega Tutorial Go web](https://github.com/bonfy/go-mega)
- [gin offical doc](https://github.com/gin-gonic/gin)
- [gin full doc](https://www.jianshu.com/p/98965b3ff638)
- [Go Gin 简明教程](https://geektutu.com/post/quick-go-gin.html)
- [gin 教程](https://youngxhui.top/categories/gin/)
- [Golang 微框架 Gin 简介](https://www.jianshu.com/p/a31e4ee25305)
- [go iris](https://wxnacy.com/2019/03/01/go-iris-simple/)

主要是这两个框架：gin / iris

略微看了一下文档，路由基本是这么配置的：

gin 的例子：

```go
func main() {
    router := gin.Default()

    // 此规则能够匹配/user/john这种格式，但不能匹配/user/ 或 /user这种格式
    router.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })

    // 但是，这个规则既能匹配/user/john/格式也能匹配/user/john/send这种格式
    // 如果没有其他路由器匹配/user/john，它将重定向到/user/john/
    router.GET("/user/:name/*action", func(c *gin.Context) {
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOK, message)
    })

    router.POST("/form_post", func(c *gin.Context) {
        message := c.PostForm("message")
        nick := c.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值

        c.JSON(200, gin.H{
            "status":  "posted",
            "message": message,
            "nick":    nick,
        })
    })

    router.Run(":8080")
}
```

iris 也差不多。

和 Node.js 的 express，Python 的 Flask 很像，不像 rails 那样是以 controller 为核心的 (Python 的 Djanjo 也是以 controller 为核心的吧?)

web 框架原理都差不多，详略，需要时再看文档。

orm 库可以用 gorm 包。

## Go 中的字符串

- [Strings, bytes, runes and characters in Go](https://blog.golang.org/strings) | [Go 语言中的字符串](https://www.jianshu.com/p/01a842787637)
- [strings — 字符串操作](https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter02/02.1.html)

总结一下就是 Go 中的字符串和 Rust 中的字符串一样，存储的是 utf-8 编码的字节码。Go 中的字符串实际就是一个 byte 切片，即 `[]byte`，这和 Rust 中字符串实际是 `Vec<u8>` 类似。

但对于单个字符类型，Go 和字符串一样仍然用 utf-8 编码，单个字符其实为 rune 类型；而 Rust 使用 unicode 编码存放。

对 Go 的字符串如何进行一个字符一个字符的遍历，使用 for...range，相当于 Rust 中的 chars() 方法。

对字符串使用 len() 方法得到的是字节个数，而不是字符个数。

如何得到真正的字符个数，使用 strings.Count() 方法，比如：

```go
fmt.Println(strings.Count("谷歌中国", "")) // 5，为实际值+1
```
