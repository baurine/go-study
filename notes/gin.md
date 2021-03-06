# Gin

- [Go gin 框架入门教程](https://www.tizi365.com/archives/244.html)

wip

## Go Web 框架

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

    // 此规则能够匹配 /user/john 这种格式，但不能匹配 /user/ 或 /user 这种格式
    router.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })

    // 但是，这个规则既能匹配 /user/john/ 格式也能匹配 /user/john/send 这种格式
    // 如果没有其他路由器匹配 /user/john，它将重定向到 /user/john/
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

## Gin 与 Middleware

## Gin 与 Jwt

## Gin 和 Swagger

## Go 和模板

- [Go 模板引擎](https://www.tizi365.com/archives/85.html)

wip

## Gin 和模板

我之前的做法：

```go
// pkg/apiserver/apiserver.go
r.LoadHTMLGlob("templates/**/*") // 模板以 .tmpl 格式放在项目根目录的 templates 目录下

// templates/sql-diagnosis/index.tmpl
{{ define "sql-diagnosis/index" }}
  <!DOCTYPE html>
  <html lang="en">
  ...
  </html>
{{ end }}

// templates/sql-diagnosis/table.tmpl
{{ define "sql-diagnosis/table" }}
  <div class="report-container">
  ...
  </div>
{{ end }}

// handler
c.HTML(http.StatusOK, "sql-diagnosis/index", tables)
```

还可以把模板作为字符常量...但有个缺点，修改了模板以后，就要重新编译，而不能像上面一样，在 debug 模式下，修改了模板只要刷新浏览器就行了。

```go
// pkg/apiserver/apiserver.go
templates := template.New("api")
r.SetHTMLTemplate(templates)  // r: Router

//...
diagnose.NewService(config, services.TiDBForwarder, services.Store).Register(endpoint, auth).RegisterTemplates(templates)

// pkg/apiserver/diagnose/diagnose.go
func (s *Service) RegisterTemplates(t *template.Template) *Service {
	_, _ = t.Parse(TemplateIndex)
	_, _ = t.Parse(TemplateTable)
	return s
}

// pkg/apiserver/diagnose/templates.go
package diagnose

const TemplateIndex = `
{{ define "sql-diagnosis/index" }}
  <!DOCTYPE html>
  <html lang="en">
  ...
  </html>
{{ end }}
`

const TemplateTable = `
{{ define "sql-diagnosis/table" }}
  <div class="report-container">
  ...
  </div>
{{ end }}
`
```

新的改法：

```go
// pkg/apiserver/apiserver.go
templates.GinLoad(r)

// templates/sqldiagnosis/index.tmpl.go
package sqldiagnosis

const Index = `
{{ define "sql-diagnosis/index" }}
  <!DOCTYPE html>
  <html lang="en">
  ...
  </html>
{{ end }}
`

// templates/sqldiagnosis/table.tmpl.go
package sqldiagnosis

const Table = `
{{ define "sql-diagnosis/table" }}
  <div class="report-container">
  ...
  </div>
{{ end }}
`

// templates/templates.go
const (
	DelimsLeft  = "{{"
	DelimsRight = "}}"
)

var DefinedTemplates = [][]string{
	{"sql-diagnosis/index", sqldiagnosis.Index},
	{"sql-diagnosis/table", sqldiagnosis.Table},
}

func GinLoad(r *gin.Engine) {
	r.Delims(DelimsLeft, DelimsRight)
	templ := template.New("").Delims(DelimsLeft, DelimsRight).Funcs(r.FuncMap)
	for _, info := range DefinedTemplates {
		name := info[0]
		text := info[1]
		t := templ.New(name)
		if _, err := t.Parse(text); err != nil {
			log.Fatal("Failed to parse template", zap.String("name", name), zap.Error(err))
		}
	}
	r.SetHTMLTemplate(templ)
}
```

后来了解到这样改的原因是为了分发 Go，Go 的一大优势就是可以把所有代码最终只编译成一个执行文件，如果用 .html 这种方式，这些文件默认是没法打包到执行文件里的。

解决办法，参考这三个 PR：

- https://github.com/pingcap-incubator/tidb-dashboard/pull/331
- https://github.com/pingcap-incubator/tidb-dashboard/pull/307
- https://github.com/pingcap-incubator/tidb-dashboard/pull/171

[Choosing A Library to Embed Static Assets in Go](https://tech.townsourced.com/post/embedding-static-files-in-go/)

思路是在 debug 模式下从文件中加载静态文件，在 release 模式下从 bundle 到 go 代码中的 assets 中加载静态文件。

http.Dir(), http.FileServer()

涉及到的新的知识点：

- Go 的条件编译

  ```go
  // a.go
  // +build dev

  // b.go
  // +build !dev
  ```

- go generate

  ```go
  // a.go
  //go:generate vfsgendev -source="github.com/pingcap-incubator/tidb-dashboard/pkg/apiserver/diagnose".Vfs
  ```

## 解析查询参数

如果查询参数是放在 url 中，比如 `/api/statements?schema_name=xxx&digest=xxx`，可以用 ShouldBindQuery(&req)，示例：

```go
type GetPlansRequest struct {
	SchemaName string `json:"schema_name" form:"schema_name"`
	Digest     string `json:"digest" form:"digest"`
	BeginTime  int    `json:"begin_time" form:"begin_time"`
	EndTime    int    `json:"end_time" form:"end_time"`
}

func (s *Service) getPlansHandler(c *gin.Context) {
	var req GetPlansRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Status(http.StatusBadRequest)
		_ = c.Error(utils.ErrInvalidRequest.WrapWithNoMessage(err))
		return
	}
```

如果是 post，可以用 ShouldBindJSON(&req)，示例：

```go
type GenerateReportRequest struct {
	StartTime        int64 `json:"start_time"`
	EndTime          int64 `json:"end_time"`
	CompareStartTime int64 `json:"compare_start_time"`
	CompareEndTime   int64 `json:"compare_end_time"`
}

func (s *Service) genReportHandler(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		_ = c.Error(apiutils.ErrInvalidRequest.WrapWithNoMessage(err))
		return
	}
```
