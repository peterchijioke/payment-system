package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"take-Home-assignment/internal/models"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var err error
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close(context.Background())

	var db *gorm.DB

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              false,
			NowFunc:                                  func() time.Time { return time.Now().UTC().Truncate(time.Microsecond) },
			Logger:                                   logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Waiting for Postgres to be ready...", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after 10 attempts: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	log.Println("Database connection established")

	log.Println("Running database migrations...")
	db.Logger = logger.Default.LogMode(logger.Silent)

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrations completed")

	log.Println("Filling missing values in existing records...")

	log.Println("Running database seeders...")
	SeedFinancialData(db)
	log.Println("Database seeders completed")

	DB = db
	log.Println("PostgreSQL connected, migrated, and seeded successfully")
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations via AutoMigrate...")
	return autoMigrate(db)
}

func autoGenerateMigration() error {
	models := []interface{}{
		models.Account{},
		models.AccountBalance{},
		models.PaymentTransaction{},
		models.LedgerEntry{},
		models.FXQuote{},
		models.WebhookEvent{},
		models.IdempotencyKey{},
	}

	var schema strings.Builder
	schema.WriteString("-- Auto-generated migration: initial_schema\n")
	schema.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
	schema.WriteString("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";\n\n")

	for _, model := range models {
		tableName := getTableName(model)
		schema.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		schema.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))

		t := reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		var fields []string
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("gorm")

			if tag == "" || strings.Contains(tag, "<-") {
				continue
			}

			fieldName := getJsonTagName(field)
			columnType := getColumnTypeFromTag(tag, field.Type.Kind())

			fields = append(fields, fmt.Sprintf("    %s %s", fieldName, columnType))
		}

		schema.WriteString(strings.Join(fields, ",\n"))
		schema.WriteString("\n);\n\n")
	}

	migrationDir := "internal/database/migrations"
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(migrationDir, "001_initial_schema.sql")
	return os.WriteFile(filename, []byte(schema.String()), 0644)
}

func getTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

func getJsonTagName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}
	return strings.ToLower(field.Name)
}

func getColumnTypeFromTag(tag string, kind reflect.Kind) string {
	if strings.Contains(tag, "type:") {
		idx := strings.Index(tag, "type:")
		endIdx := strings.Index(tag[idx:], ";")
		if endIdx == -1 {
			endIdx = len(tag[idx:])
		}
		return tag[idx+5 : idx+endIdx]
	}

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

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
		&models.AccountBalance{},
		&models.PaymentTransaction{},
		&models.LedgerEntry{},
		&models.FXQuote{},
		&models.WebhookEvent{},
		&models.IdempotencyKey{},
	)
}
