# Gin

- [Go gin 框架入门教程](https://www.tizi365.com/archives/244.html)

wip

## Gin 与 Jwt

## Gin 和 Swagger

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
