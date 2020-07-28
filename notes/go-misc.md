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

## 变量分配在栈还是堆

c/c++ 中，用 malloc/new 在堆上分配空间，其余在栈上分配。Go 不太一样，new 的对象也有可能在栈上分配，而不是 new 的对象有可能在堆上分配，由 Go 编译器决定，它会进行逃逸分析，简单地说就是如果变量仅在函数内作用，不会被引用到函数外，那么就在栈上分配，不管是 new 还不是 new 出来的，反之，则在堆上分配。

所以 Go 让你完全不用再考虑栈还是堆的问题...

## Go Context

- [6.1 上下文 Context](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/)

wip

## Go 时间解析

Go 的时间解析有点奇芭，和其它语言很不一样。它不使用 "YYYY-MM-DD HH:mm:ss" 这样的模板，而是使用了一个特定时间作为模板 (layout)，即 2006 年 1 月 2 日，下午 3 时 4 分 5 秒。比如：

```go
fmt.Println(time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")) //2014-01-07 09:32:12

dateStr := "2016-07-14 14:24:51"
timestamp1, _ := time.Parse("2006-01-02 15:04:05", dateStr)
timestamp2, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
fmt.Println(timestamp1, timestamp2)               //2016-07-14 14:24:51 +0000 UTC 2016-07-14 14:24:51 +0800 CST
fmt.Println(timestamp1.Unix(), timestamp2.Unix()) //1468506291 1468477491

p := fmt.Println
t := time.Now()
p(t.Format("2006-01-02T15:04:05Z07:00"))
p(t.Format("3:04PM"))
p(t.Format("Mon Jan _2 15:04:05 2006"))
p(t.Format("2006-01-02T15:04:05.999999-07:00"))
```

将在 url 作为查询参数的时间戳转换成人类易读的时间格式：

```go
startTimeStr := c.Query("start_time")
tsSec, err := strconv.ParseInt(startTimeStr, 10, 64)
if err != nil {
  _ = c.Error(err)
  return
}
startTime := time.Unix(tsSec, 0)
```

## Go 文件操作

- [Go 文件操作大全 (译)](https://colobu.com/2016/10/12/go-file-operations/)

主要涉及 os | io | io/ioutil | bufio 包。

## Host 在 HTTP Request Header 中的作用

(其实这个内容与 Go 无关，是属于基本的网络知识)

Related:

- https://github.com/pingcap/telemetry-log-collector/pull/1
- https://github.com/yeqown/fasthttp-reverse-proxy/pull/11
- [HTTP 协议 Host 请求头的作用](https://blog.csdn.net/codejas/article/details/82844032)

> Host 是 HTTP 1.1 协议中新增的一个请求头，主要用来实现虚拟主机技术。

> 虚拟主机（virtual hosting）即共享主机（shared web hosting），可以利用虚拟技术把一台完整的服务器分成若干个主机，因此可以在单一主机上运行多个网站或服务。

> 举个栗子，有一台 ip 地址为 61.135.169.125 的服务器，在这台服务器上部署着谷歌、百度、淘宝的网站。为什么我们访问 https://www.google.com 时，看到的是 Google 的首页而不是百度或者淘宝的首页？原因就是 Host 请求头决定着访问哪个虚拟主机。

(又有点疑惑了，对于 fasthttp-reverse-proxy，不设置 host 也能工作呢...有待进一步研究)

TODO: add code example
