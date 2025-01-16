package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lastnight2022-cc/tools_lib/generator/xorm"
)

func main() {
	dsn := flag.String("dsn", "", "Data Source Name for database connection")
	output := flag.String("output", "./models", "Output directory for generated files")
	tableName := flag.String("table", "", "Specific table name to generate (optional)")
	flag.Parse()

	if *dsn == "" {
		fmt.Println("Error: DSN is required.")
		os.Exit(1)
	}

	err := generator.GenerateStructs(*dsn, *output, *tableName)
	if err != nil {
		fmt.Printf("Error generating structs: %v\n", err)
		os.Exit(1)
	}
}
