package xorm

import "text/template"

var structTemplate = template.Must(template.New("struct").Parse(`
type {{.Name}} struct {
{{range .Columns}}
    {{.FieldName}} {{.FieldType}} ` + "`xorm:\"{{.ColumnName}}\"`" + `
{{end}}
}
`))

// 可以添加更多的模板变量和逻辑来丰富生成的内容
