package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"take-Home-assignment/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generate_migration.go <migration_name>")
	}

	migrationName := os.Args[1]
	filename := fmt.Sprintf("internal/database/migrations/%s_%s.sql", time.Now().Format("20060102150405"), migrationName)

	var schema strings.Builder

	schema.WriteString("-- Auto-generated migration: " + migrationName + "\n")
	schema.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
	schema.WriteString("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";\n\n")

	models := []interface{}{
		models.Account{},
		models.AccountBalance{},
		models.PaymentTransaction{},
		models.LedgerEntry{},
		models.FXQuote{},
		models.WebhookEvent{},
		models.IdempotencyKey{},
	}

	for _, model := range models {
		tableName := getTableName(model)
		schema.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		schema.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))

		v := reflect.ValueOf(model)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		var fields []string
		for i := 0; i < v.Type().NumField(); i++ {
			field := v.Type().Field(i)
			tag := field.Tag.Get("gorm")
			jsonTag := field.Tag.Get("json")

			if tag == "" || strings.Contains(tag, "<-") {
				continue
			}

			columnType := getColumnType(field.Type.Kind())
			if strings.Contains(tag, "type:") {
				if idx := strings.Index(tag, "type:"); idx != -1 {
					endIdx := strings.Index(tag[idx:], ";")
					if endIdx == -1 {
						endIdx = len(tag[idx:])
					}
					columnType = tag[idx+5 : idx+endIdx]
				}
			}

			fieldName := getFieldName(field, jsonTag)
			fields = append(fields, fmt.Sprintf("    %s %s", fieldName, columnType))
		}

		schema.WriteString(strings.Join(fields, ",\n"))
		schema.WriteString("\n);\n\n")
	}

	err := os.WriteFile(filename, []byte(schema.String()), 0644)
	if err != nil {
		log.Fatalf("Failed to write migration file: %v", err)
	}

	fmt.Printf("Migration file created: %s\n", filename)
}

func getTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func getFieldName(field reflect.StructField, jsonTag string) string {
	if jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}
	return strings.ToLower(field.Name)
}

func getColumnType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "DECIMAL(38, 12)"
	case reflect.Bool:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}
