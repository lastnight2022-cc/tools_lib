package xorm

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
	"xorm.io/xorm/schemas"

	_ "github.com/go-sql-driver/mysql" // 确保导入了合适的驱动
	"xorm.io/xorm"
)

// GenerateStructs connects to the database and generates Go structs based on tables.
func GenerateStructs(dsn, outputDir, tableName string) error {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer engine.Close()

	// 获取表列表或单个表信息
	tables, err := getTables(engine, tableName)
	if err != nil {
		return err
	}

	// 遍历每个表并生成相应的 Go struct 文件
	for _, table := range tables {
		content, err := generateTableStruct(engine, table)
		if err != nil {
			return err
		}

		outputPath := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(table))
		if err := ioutil.WriteFile(outputPath, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("Generated file: %s\n", outputPath)
	}

	return nil
}

// getTables retrieves a list of table names from the database.
func getTables(engine *xorm.Engine, specificTable string) ([]string, error) {
	var tables []string
	if specificTable != "" {
		tables = append(tables, specificTable)
	} else {
		// 这里需要根据具体的数据库类型调整 SQL 查询语句
		sql := "SHOW TABLES"
		rows, err := engine.Query(sql)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			for _, col := range row {
				tables = append(tables, string(col))
			}
		}
	}
	return tables, nil
}

// generateTableStruct generates a Go struct for a given table.
func generateTableStruct(engine *xorm.Engine, table string) (string, error) {
	// 使用 engine.DatabaseMetas() 获取表的元数据
	dbMetas, err := engine.DBMetas()
	if err != nil {
		return "", err
	}

	var columns []*schemas.Column
	for _, meta := range dbMetas {
		if meta.Name == table {
			columns = meta.Columns()
			break
		}
	}

	if len(columns) == 0 {
		return "", fmt.Errorf("table %s not found in database", table)
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(table)))

	for _, column := range columns {
		builder.WriteString(fmt.Sprintf("    %s %s `xorm:\"%s\"`\n",
			capitalize(column.Name),
			mapDbTypeToGoType(column.SQLType.Name),
			column.Name,
		))
	}
	builder.WriteString("}\n")

	return builder.String(), nil
}

// capitalize capitalizes the first letter of a string.
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

// mapDbTypeToGoType maps database column types to Go types.
func mapDbTypeToGoType(dbType string) string {
	switch dbType {
	case "int", "tinyint", "smallint", "mediumint", "bigint":
		return "int64"
	case "float", "double":
		return "float64"
	case "char", "varchar", "text", "longtext":
		return "string"
	case "datetime", "timestamp":
		return "time.Time"
	default:
		return "interface{}"
	}
}
