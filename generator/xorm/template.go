package xorm

import "text/template"

// 修改: 添加了对包名和导入语句的支持
var structTemplate = template.Must(template.New("struct").Parse(`
package {{.PackageName}}
{{if .Imports}}
import (
{{range .Imports}}
    "{{.}}"
{{end}}
)
{{end}}
type {{.Name}} struct {
{{range .Columns}}
    {{.FieldName}} {{.FieldType}} ` + "`xorm:\"{{.ColumnName}}\"`" + ` // {{.ColumnName}} {{.ColumnType}}
{{end}}
}
`))
