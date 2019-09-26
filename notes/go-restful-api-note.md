# Build Go RESTful API by gin Note

掘金小册[《基于 Go 语言构建企业级的 RESTful API 服务》](https://juejin.im/book/5b0778756fb9a07aa632301e)的笔记，相应的实践代码在 code/apiserver 目录。

这个小册虽然写得不是很细心，但还是比较全面地覆盖了一个完整的 web 开发所涉及的大部分流程。

详细内容：

- 准备阶段

  - 如何安装和配置 Go 开发环境
  - 如何安装和配置 Vim IDE

- 设计阶段

  - API 构建技术选型
  - API 基本原理
  - API 规范设计

- 开发阶段

  - 如何读取配置文件 - 使用 viper 库
  - 如何管理和记录日志 - 使用 log 库
  - 如何做数据库的 CURD 操作 - 使用 gorm 库
  - 如何自定义错误 Code
  - 如何读取和返回 HTTP 请求
  - 如何进行业务逻辑开发
  - 如何对请求插入自己的处理逻辑 - middleware
  - 如何进行 API 身份验证 - jwt
  - 如何进行 HTTPS 加密 - https
  - 如何用 Makefile 管理编译
  - 如何给 API 命令添加版本功能
  - 如何管理 API 命令 - admin.sh
  - 如何生成 Swagger 在线文档

- 测试阶段

  - 如何进行单元测试
  - 如何进行性能测试（函数性能）
  - 如何做性能分析
  - API 性能测试和调优

- 部署阶段
  - 如何用 Nginx 部署 API 服务
  - 如何做 API 高可用

## RESTful API 介绍

略。REST vs RPC。

## API 流程和代码结构

略。

一般流程：

- (配置连接数据库)
- (配置 log)
- 生成服务器对象：`g := gin.New()`，node.js 的 `var app = express()`
- 注册中间件
- 注册路由
- 启动服务器，监听端口

## Go API 开发环境配置

略。

## 基础 1：启动一个最简单的 RESTful API 服务器

REST Web 框架选择，选择使用 gin。

加载路由，支持 group。

```go
"apiserver/handler/sd"

....

// The health check handlers
svcd := g.Group("/sd")
{
    svcd.GET("/health", sd.HealthCheck)
    svcd.GET("/disk", sd.DiskCheck)
    svcd.GET("/cpu", sd.CPUCheck)
    svcd.GET("/ram", sd.RAMCheck)
}
```

使用 `g.Use()` 配置中间件：

```go
g.Use(gin.Recovery())
g.Use(middleware.NoCache)
g.Use(middleware.Options)
g.Use(middleware.Secure)
```

API 服务器健康状态自检，启动一个 goroutine 去访问自己的 api (使用 `http.GET()` 方法)，看是否成功。

```go
func main() {
    ....

    // Ping the server to make sure the router is working.
    go func() {
        if err := pingServer(); err != nil {
            log.Fatal("The router has no response, or it might took too long to start up.", err)
        }
        log.Print("The router has been deployed successfully.")
    }()
    ....
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
    for i := 0; i < 10; i++ {
        // Ping the server by sending a GET request to `/health`.
        resp, err := http.Get("http://127.0.0.1:8080" + "/sd/health")
        if err == nil && resp.StatusCode == 200 {
            return nil
        }

        // Sleep for a second to continue the next ping.
        log.Print("Waiting for the router, retry in 1 second.")
        time.Sleep(time.Second)
    }
    return errors.New("Cannot connect to the router.")
}
```

cURL 工具测试 API。

```
-X/--request [GET|POST|PUT|DELETE|...]  指定请求的 HTTP 方法
-H/--header                           指定请求的 HTTP Header
-d/--data                             指定请求的 HTTP 消息体（Body）
-v/--verbose                          输出详细的返回信息
-u/--user                             指定账号、密码
-b/--cookie                           读取 cookie
```

一个示例：

```shell
$ curl -v -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/user -d'{"username":"admin","password":"admin1234"}'
```

可以用 jq 命令 (需要安装 `brew install jq`) 将得到的 json 在命令行进行格式化。

```shell
$ curl ... | jq .
```

## 基础 2：配置文件读取

使用了 viper 库，详细使用看文档。

## 基础 3：记录和管理 API 日志

使用了作者自己的 lexkong/log 库。详略。

## 基础 4：安装 MySQL 并初始化表

略。见 MySQL 笔记。

## 基础 5：初始化 MySQL 数据库并建立连接

使用了 gorm 库，详略。

## 基础 6：自定义业务错误信息

定义了很多 errno。详略。

## 基础 7：读取和返回 HTTP 请求

gin 读取 request 中的参数提供的一些方法：

- c.Param(key) - 返回 path 中的参数值
- c.Query(key) - 返回 url 中 query 中的值
- c.DefaultKey(key, defVal) - 相比 Query() 如果 key 不存在返回默认值
- c.Bind() - 检查 Content-Type 类型，将消息体作为指定的格式解析到 Go struct 变量中
- c.GetHeader() - 获取 Http Header

(rails 则是全部统一成了 params 方法)

c 是 `*gin.Context` 类型。

gin 返回 response 的一些方法：

- c.String(statusCode, body)
- c.JSON(statusCode, body)
- ...

如果是做 API，那毫无疑问几乎就是用 c.JSON() 方法了。本小册将发送 response 的逻辑封装成了 `SendResponse(c *gin.Context, err error, data interface{})` 方法。

```go
// handler/handler.go
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
```

## 基础 8：用户业务逻辑处理

对 user model 进行 CRUD。详略。

路由：

```go
// 用户路由设置
u := g.Group("/v1/user")
{
  u.POST("", user.Create)         // 创建用户
  u.DELETE("/:id", user.Delete)   // 删除用户
  u.PUT("/:id", user.Update)      // 更新用户
  u.GET("", user.List)            // 用户列表
  u.GET("/:username", user.Get)   // 获取指定用户的详细信息
}
```

## 基础 9：HTTP 调用添加自定义处理逻辑

添加中间件，加了给 Header 插件 RequestId 和 log 每个请求 body 的中间件。详略。

```go
func main() {
    ...
    // Routes.
    router.Load(
        // Cores.
        g,

        // Middlwares.
        middleware.RequestId(),
        middleware.Logging()
    )
    ...
}
```

## 基础 10：API 身份验证

使用 jwt 作为身份验证。在 "/login" 请求中生成 jwt token，在 "/v1/user" 请求中验证 jwt token。详略。

```go
g.POST("/login", user.Login)

u := g.Group("/v1/user")
u.Use(middleware.Auth())
{
  u.POST("", user.Create)
  u.DELETE("/:id", user.Delete)
  u.PUT("/:id", user.Update)
  u.GET("", user.List)
  u.GET("/:username", user.Get)
}
```

可以把中间件仅作用于部分请求。

## 进阶 1：用 HTTPS 加密 API 请求

net/http 包提供了 ListenAndServerTLS() 方法创建 https 服务端，和 http 的区别是提供数据证书 certFile 和密钥文件 keyFile。

```go
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error
```

测试时可以用 openssl 命令生成密钥和自签名的证书。

```shell
$ penssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"
```

代码：

```go
cert := viper.GetString("tls.cert")
key := viper.GetString("tls.key")
if cert != "" && key != "" {
  go func() {
    log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.port"))
    log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
  }()
}
```

使用 curl 测试时，使用 `--cacert` `--cert` `--key` 参数声明 certFile 和 keyFile。

```shell
curl -XGET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ" -H "Content-Type: application/json" https://127.0.0.1:8081/v1/user --cacert conf/server.crt --cert conf/server.crt --key conf/server.key | jq .
```

## 进阶 2：用 Makefile 管理 API 项目

编写 Makefile，详略。把笔记补充到 programming-basic repo 中。

## 进阶 3：给 API 命令增加版本功能

增加 `-v` 命令选择，输出该程序的版本等信息。这些版本是在执行 `go build` 进行编译时通过 `-ldflags -X importpath.name=value` 选项直接写进二进制包中的。

详略。

## 进阶 4：给 API 增加启动脚本

编写 admin.sh 来启动/停止 apiserver。详略。

## 进阶 5：基于 Nginx 的 API 部署方案

使用 nginx 进行反向代理及负载均衡。部分 nginx 配置如下：

```nginx
http {
  # ...

  upstream apiserver {
      server 127.0.0.1:8080;
      server 127.0.0.1:8082;
  }

  server {
      listen      80;
      server_name  localhost;

      location / {
          proxy_pass  http://apiserver;
      }
  }
}
```

## 进阶 6：API 高可用方案

防止 nginx 的单点故障，布署多个 nginx 节点。

## 进阶 7：go test 测试你的代码

为 util 包中的函数编写单元测试及 benchmark。查看性能及生成函数调用图，生成 svg (需要安装 graphviz)，生成测试覆盖率。详略。

## 进阶 8：API 性能分析

使用 pprof 进行 api 性能分析，详略，需要时再看。

## 进阶 9：生成 Swagger 在线文档

集成 swagger，详略。

## 进阶 10：API 性能测试和调优

使用 wrk 工具，对 web 进行并发数和 QPS 的压测。

## 拓展 1：Go 开发技巧

略。

## 拓展 2：Go 规范指南

略。

最后，总结一下，gin 框架总体和 node.js 的 express，python 的 flask 都比较相似，比较轻量级，以 handler 和 model 为核心。而 rails 以 controller 为核心，rails 的 ActiveRecord 相比其它 ORM 强大和方便很多。
