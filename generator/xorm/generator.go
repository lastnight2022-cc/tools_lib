package xorm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	// 检查并创建输出目录
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// 获取表列表或单个表信息
	tables, err := getTables(engine, tableName)
	if err != nil {
		return err
	}

	// 遍历每个表并生成相应的 Go struct 文件
	for _, table := range tables {
		// 使用目录名作为包名
		packageName := strings.ToLower(filepath.Base(outputDir))
		content, err := generateTableStruct(engine, table, packageName)
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
func generateTableStruct(engine *xorm.Engine, table string, packageName string) (string, error) {
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

	var imports []string
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName)) // 增加包名

	// 添加导入语句
	for _, column := range columns {
		goType := mapDbTypeToGoType(column.SQLType.Name)
		if goType == "time.Time" {
			if !sliceContains(imports, "time") {
				imports = append(imports, "time")
			}
		}
	}
	if len(imports) > 0 {
		builder.WriteString("import (\n")
		for _, imp := range imports {
			builder.WriteString(fmt.Sprintf("    \"%s\"\n", imp))
		}
		builder.WriteString(")\n\n")
	}

	builder.WriteString(fmt.Sprintf("type %s struct {\n", capitalize(table)))

	for _, column := range columns {
		var tag string
		if column.IsPrimaryKey && column.IsAutoIncrement {
			tag = "`xorm:\"pk autoincr\"`"
		} else {
			tag = fmt.Sprintf("`xorm:\"%s", column.Name)
			if column.Nullable {
				tag += " omitempty"
			}
			tag += "\"`"
		}
		builder.WriteString(fmt.Sprintf("    %s %s %s // %s %s\n",
			capitalize(column.Name),
			mapDbTypeToGoType(column.SQLType.Name),
			tag,
			column.Name,
			column.SQLType.Name,
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
	// 将 dbType 转换为小写
	dbType = strings.ToLower(dbType)
	// 判断是否为 unsigned 类型
	isUnsigned := strings.Contains(dbType, "unsigned")
	if isUnsigned {
		dbType = strings.Replace(dbType, "unsigned", "", -1)
		dbType = strings.TrimSpace(dbType)
	}
	switch dbType {
	case "int":
		if isUnsigned {
			return "uint"
		}
		return "int64"
	case "tinyint":
		if isUnsigned {
			return "uint8"
		}
		return "int64"
	case "smallint":
		if isUnsigned {
			return "uint16"
		}
		return "int64"
	case "mediumint":
		if isUnsigned {
			return "uint32"
		}
		return "int64"
	case "bigint":
		if isUnsigned {
			return "uint64"
		}
		return "int64"
	case "float", "double", "decimal":
		return "float64"
	case "char", "varchar", "text", "longtext", "tinytext", "mediumtext":
		return "string"
	case "datetime", "timestamp", "date", "time":
		return "time.Time"
	case "boolean", "tinyint(1)":
		return "bool"
	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "[]byte"
	default:
		return "interface{}"
	}
}

// 添加辅助函数 sliceContains
func sliceContains(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}
