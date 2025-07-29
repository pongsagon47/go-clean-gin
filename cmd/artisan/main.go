// cmd/artisan/main.go - Complete Laravel-style CLI tool
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/pkg/database"
	"go-clean-gin/pkg/logger"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	action = flag.String("action", "", "Action: make:migration, make:seeder, make:model, make:package, migrate, migrate:rollback, migrate:status")
	name   = flag.String("name", "", "Migration/Seeder/Model/Package name")
	table  = flag.String("table", "", "Table name for migration")
	create = flag.Bool("create", false, "Create table migration")
	fields = flag.String("fields", "", "Fields for migration (name:type,email:string)")
	count  = flag.Int("count", 1, "Number of migrations to rollback")
	help   = flag.Bool("help", false, "Show help")
)

func main() {
	flag.Parse()

	if *help || *action == "" {
		showHelp()
		return
	}

	switch *action {
	case "make:migration":
		if *name == "" {
			fmt.Println("❌ Migration name is required")
			fmt.Println("Usage: go run cmd/artisan/main.go -action=make:migration -name=migration_name")
			os.Exit(1)
		}
		createMigration(*name, *table, *create, *fields)

	case "make:seeder":
		if *name == "" {
			fmt.Println("❌ Seeder name is required")
			fmt.Println("Usage: go run cmd/artisan/main.go -action=make:seeder -name=seeder_name")
			os.Exit(1)
		}
		createSeeder(*name, *table)

	case "make:model":
		if *name == "" {
			fmt.Println("❌ Model name is required")
			fmt.Println("Usage: go run cmd/artisan/main.go -action=make:model -name=model_name")
			os.Exit(1)
		}
		createModel(*name, *fields)

	case "make:package":
		if *name == "" {
			fmt.Println("❌ Package name is required")
			fmt.Println("Usage: go run cmd/artisan/main.go -action=make:package -name=package_name")
			os.Exit(1)
		}
		createPackage(*name)

	case "migrate":
		runMigrations()

	case "migrate:rollback":
		rollbackMigrations(*count)

	case "migrate:status":
		showMigrationStatus()

	case "db:seed":
		runSeeders(*name)

	default:
		fmt.Printf("❌ Unknown action: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
}

func createMigration(migrationName, tableName string, isCreate bool, fieldList string) {
	timestamp := time.Now().Format("2006_01_02_150405")
	fileName := fmt.Sprintf("%s_%s.go", timestamp, toSnakeCase(migrationName))

	// Create migrations directory if not exists
	migrationsDir := "internal/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create migrations directory: %v\n", err)
		os.Exit(1)
	}

	filePath := filepath.Join(migrationsDir, fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("❌ Migration file already exists: %s\n", filePath)
		os.Exit(1)
	}

	// Parse fields
	var parsedFields []Field
	if fieldList != "" {
		fieldPairs := strings.Split(fieldList, ",")
		for _, pair := range fieldPairs {
			parts := strings.Split(strings.TrimSpace(pair), ":")
			if len(parts) == 2 {
				parsedFields = append(parsedFields, Field{
					Name: strings.TrimSpace(parts[0]),
					Type: strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	// Create migration data
	data := MigrationData{
		ClassName:   toPascalCase(migrationName),
		TableName:   tableName,
		Timestamp:   timestamp,
		Description: migrationName,
		Fields:      parsedFields,
		Version:     fmt.Sprintf("%s_%s", timestamp, migrationName),
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("❌ Failed to create migration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Choose template
	var tmpl *template.Template
	if isCreate && tableName != "" {
		tmpl = template.Must(template.New("create_table").Funcs(templateFuncs).Parse(createTableTemplate))
	} else if tableName != "" {
		tmpl = template.Must(template.New("alter_table").Funcs(templateFuncs).Parse(alterTableTemplate))
	} else {
		tmpl = template.Must(template.New("migration").Funcs(templateFuncs).Parse(migrationTemplate))
	}

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("❌ Failed to generate migration file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Migration created: %s\n", filePath)
	fmt.Printf("📝 Class: %s\n", data.ClassName)
	if tableName != "" {
		fmt.Printf("🗂️  Table: %s\n", tableName)
	}
}

func createSeeder(seederName, tableName string) {
	if !strings.HasSuffix(seederName, "Seeder") {
		seederName += "Seeder"
	}

	fmt.Println(seederName)

	fileName := fmt.Sprintf("%s.go", toSnakeCase(seederName))
	// fileName := fmt.Sprintf("%s.go", toSnakeCase(strings.TrimSuffix(seederName, "Seeder")))

	// Create seeders directory if not exists
	seedersDir := "internal/seeders"
	if err := os.MkdirAll(seedersDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create seeders directory: %v\n", err)
		os.Exit(1)
	}

	filePath := filepath.Join(seedersDir, fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("❌ Seeder file already exists: %s\n", filePath)
		os.Exit(1)
	}

	data := SeederData{
		ClassName: seederName,
		TableName: tableName,
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("❌ Failed to create seeder file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute template
	tmpl := template.Must(template.New("seeder").Parse(seederTemplate))
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("❌ Failed to generate seeder file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Seeder created: %s\n", filePath)
	fmt.Printf("📝 Class: %s\n", data.ClassName)
	if tableName != "" {
		fmt.Printf("🗂️  Table: %s\n", tableName)
	}
}

func createModel(modelName, fieldList string) {
	// Generate entity struct name
	entityName := toPascalCase(modelName)
	tableName := strings.ToLower(toSnakeCase(entityName)) + "s" // posts, users, etc.

	fileName := fmt.Sprintf("%s.go", strings.ToLower(entityName))

	// Create entity directory if not exists
	entityDir := "internal/entity"
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create entity directory: %v\n", err)
		os.Exit(1)
	}

	filePath := filepath.Join(entityDir, fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("❌ Entity file already exists: %s\n", filePath)
		os.Exit(1)
	}

	// Parse fields
	var parsedFields []Field
	if fieldList != "" {
		fieldPairs := strings.Split(fieldList, ",")
		for _, pair := range fieldPairs {
			parts := strings.Split(strings.TrimSpace(pair), ":")
			if len(parts) == 2 {
				parsedFields = append(parsedFields, Field{
					Name: strings.TrimSpace(parts[0]),
					Type: strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	// Create entity data
	data := EntityData{
		EntityName: entityName,
		TableName:  tableName,
		Fields:     parsedFields,
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("❌ Failed to create entity file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute template
	tmpl := template.Must(template.New("entity").Funcs(templateFuncs).Parse(entityTemplate))
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("❌ Failed to generate entity file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Entity created: %s\n", filePath)
	fmt.Printf("📝 Entity: %s\n", entityName)
	fmt.Printf("🗂️  Table: %s\n", tableName)
}

func createPackage(packageName string) {
	// Convert to lowercase for package name
	pkgName := strings.ToLower(packageName)
	entityName := toPascalCase(packageName)

	// Create package directory
	packageDir := filepath.Join("internal", pkgName)
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create package directory: %v\n", err)
		os.Exit(1)
	}

	// Check if package already exists
	files := []string{"handler.go", "port.go", "repository.go", "usecase.go"}
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(packageDir, file)); err == nil {
			fmt.Printf("❌ Package '%s' already exists (found %s)\n", pkgName, file)
			os.Exit(1)
		}
	}

	packageData := PackageData{
		PackageName: pkgName,
		EntityName:  entityName,
	}

	// Create handler.go
	if err := createFileFromTemplate(
		filepath.Join(packageDir, "handler.go"),
		handlerTemplate,
		packageData,
	); err != nil {
		fmt.Printf("❌ Failed to create handler.go: %v\n", err)
		os.Exit(1)
	}

	// Create port.go
	if err := createFileFromTemplate(
		filepath.Join(packageDir, "port.go"),
		portTemplate,
		packageData,
	); err != nil {
		fmt.Printf("❌ Failed to create port.go: %v\n", err)
		os.Exit(1)
	}

	// Create repository.go
	if err := createFileFromTemplate(
		filepath.Join(packageDir, "repository.go"),
		repositoryTemplate,
		packageData,
	); err != nil {
		fmt.Printf("❌ Failed to create repository.go: %v\n", err)
		os.Exit(1)
	}

	// Create usecase.go
	if err := createFileFromTemplate(
		filepath.Join(packageDir, "usecase.go"),
		usecaseTemplate,
		packageData,
	); err != nil {
		fmt.Printf("❌ Failed to create usecase.go: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Package created: internal/%s/\n", pkgName)
	fmt.Printf("📁 Files created:\n")
	fmt.Printf("  - internal/%s/handler.go\n", pkgName)
	fmt.Printf("  - internal/%s/port.go\n", pkgName)
	fmt.Printf("  - internal/%s/repository.go\n", pkgName)
	fmt.Printf("  - internal/%s/usecase.go\n", pkgName)
	fmt.Printf("🎯 Entity: %s\n", entityName)
}

func createFileFromTemplate(filePath, templateContent string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("template").Funcs(templateFuncs).Parse(templateContent))
	return tmpl.Execute(file, data)
}

func runMigrations() {
	fmt.Println("⬆️  Running migrations...")

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("❌ Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		fmt.Printf("❌ Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Migrations completed successfully")
}

func rollbackMigrations(count int) {
	fmt.Printf("⬇️  Rolling back %d migration(s)...\n", count)

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("❌ Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Rollback migrations
	if err := database.RollbackMigrations(db, count); err != nil {
		fmt.Printf("❌ Rollback failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Rollback completed successfully")
}

func showMigrationStatus() {
	fmt.Println("📊 Checking migration status...")

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("❌ Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Show migration status
	if err := database.GetMigrationStatus(db); err != nil {
		fmt.Printf("❌ Failed to get migration status: %v\n", err)
		os.Exit(1)
	}
}

func runSeeders(seederName string) {
	fmt.Println("🌱 Running seeders...")

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("❌ Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Run seeders
	if err := database.SeedData(db, seederName); err != nil {
		fmt.Printf("❌ Seeding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Seeding completed successfully")
}

func showHelp() {
	fmt.Println("🎨 Go Clean Gin - Artisan CLI (Laravel Style)")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/artisan/main.go -action=<action> [options]")
	fmt.Println("")
	fmt.Println("Available Actions:")
	fmt.Println("  make:migration     Create a new migration file")
	fmt.Println("  make:seeder        Create a new seeder file")
	fmt.Println("  make:model         Create a new entity model file")
	fmt.Println("  make:package       Create a new package with handler, usecase, repository, port")
	fmt.Println("  migrate            Run pending migrations")
	fmt.Println("  migrate:rollback   Rollback migrations")
	fmt.Println("  migrate:status     Show migration status")
	fmt.Println("  db:seed            Run database seeders")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -name string       Migration/Seeder/Model/Package name")
	fmt.Println("  -table string      Table name")
	fmt.Println("  -create            Create table migration")
	fmt.Println("  -fields string     Fields (name:string,email:string)")
	fmt.Println("  -count int         Number of migrations to rollback (default: 1)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Create table migration")
	fmt.Println("  go run cmd/artisan/main.go -action=make:migration -name=create_users_table -create -table=users -fields=\"name:string,email:string\"")
	fmt.Println("")
	fmt.Println("  # Create entity model")
	fmt.Println("  go run cmd/artisan/main.go -action=make:model -name=User -fields=\"name:string,email:string,age:int\"")
	fmt.Println("")
	fmt.Println("  # Create package (handler, usecase, repository, port)")
	fmt.Println("  go run cmd/artisan/main.go -action=make:package -name=Product")
	fmt.Println("")
	fmt.Println("  # Add column migration")
	fmt.Println("  go run cmd/artisan/main.go -action=make:migration -name=add_phone_to_users -table=users -fields=\"phone:string\"")
	fmt.Println("")
	fmt.Println("  # Run migrations")
	fmt.Println("  go run cmd/artisan/main.go -action=migrate")
	fmt.Println("")
	fmt.Println("  # Rollback last 2 migrations")
	fmt.Println("  go run cmd/artisan/main.go -action=migrate:rollback -count=2")
	fmt.Println("")
	fmt.Println("  # Create seeder")
	fmt.Println("  go run cmd/artisan/main.go -action=make:seeder -name=UserSeeder -table=users")
}

// Helper types and functions
type MigrationData struct {
	ClassName   string
	TableName   string
	Timestamp   string
	Description string
	Fields      []Field
	Version     string
}

type Field struct {
	Name string
	Type string
}

type SeederData struct {
	ClassName string
	TableName string
}

type EntityData struct {
	EntityName string
	TableName  string
	Fields     []Field
}

type PackageData struct {
	PackageName string
	EntityName  string
}

// Template functions
var templateFuncs = template.FuncMap{
	"toSQLType":        toSQLType,
	"toGoType":         toGoType,
	"toPascalCase":     toPascalCase,
	"getGormTag":       getGormTag,
	"getValidationTag": getValidationTag,
	"hasDecimalField":  hasDecimalField,
}

func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(strings.ToLower(word))
	}

	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toSQLType(goType string) string {
	switch strings.ToLower(goType) {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "int", "integer":
		return "INTEGER"
	case "int64", "bigint":
		return "BIGINT"
	case "float", "float64":
		return "DOUBLE PRECISION"
	case "decimal":
		return "DECIMAL(10,2)"
	case "bool", "boolean":
		return "BOOLEAN"
	case "uuid":
		return "UUID"
	case "timestamp", "time":
		return "TIMESTAMP WITH TIME ZONE"
	case "date":
		return "DATE"
	case "json", "jsonb":
		return "JSONB"
	default:
		return "VARCHAR(255)"
	}
}

func toGoType(fieldType string) string {
	switch strings.ToLower(fieldType) {
	case "string":
		return "string"
	case "text":
		return "string"
	case "int", "integer":
		return "int"
	case "int64", "bigint":
		return "int64"
	case "float", "float64":
		return "float64"
	case "decimal":
		return "decimal.Decimal"
	case "bool", "boolean":
		return "bool"
	case "uuid":
		return "uuid.UUID"
	case "timestamp", "time":
		return "time.Time"
	case "date":
		return "time.Time"
	case "json", "jsonb":
		return "map[string]interface{}"
	default:
		return "string"
	}
}

func getGormTag(fieldType string) string {
	switch strings.ToLower(fieldType) {
	case "string":
		return "not null"
	case "text":
		return "type:text"
	case "int", "integer":
		return "not null"
	case "int64", "bigint":
		return "type:bigint;not null"
	case "float", "float64":
		return "type:double precision;not null"
	case "decimal":
		return "type:decimal(10,2);not null"
	case "bool", "boolean":
		return "default:false"
	case "uuid":
		return "type:uuid;not null"
	case "timestamp", "time":
		return "type:timestamp with time zone"
	case "date":
		return "type:date"
	case "json", "jsonb":
		return "type:jsonb;default:'{}'"
	default:
		return "not null"
	}
}

func getValidationTag(fieldType string) string {
	switch strings.ToLower(fieldType) {
	case "string":
		return "required,min=1,max=255"
	case "text":
		return "required"
	case "int", "integer":
		return "required,min=0"
	case "int64", "bigint":
		return "required,min=0"
	case "float", "float64":
		return "required,min=0"
	case "decimal":
		return "required,min=0"
	case "bool", "boolean":
		return ""
	case "uuid":
		return "required"
	case "timestamp", "time":
		return ""
	case "date":
		return ""
	case "json", "jsonb":
		return ""
	default:
		return "required"
	}
}

func hasDecimalField(fields []Field) bool {
	for _, field := range fields {
		if strings.ToLower(field.Type) == "decimal" {
			return true
		}
	}
	return false
}

// Templates
const migrationTemplate = `package migrations

import (
	"gorm.io/gorm"
)

// {{.ClassName}} migration
type {{.ClassName}} struct{}

// Up runs the migration
func (m *{{.ClassName}}) Up(db *gorm.DB) error {
	// TODO: Implement your migration logic here
	return nil
}

// Down rolls back the migration  
func (m *{{.ClassName}}) Down(db *gorm.DB) error {
	// TODO: Implement your rollback logic here
	return nil
}

// Description returns migration description
func (m *{{.ClassName}}) Description() string {
	return "{{.Description}}"
}

// Version returns migration version
func (m *{{.ClassName}}) Version() string {
	return "{{.Version}}"
}

// Auto-register migration
func init() {
	Register(&{{.ClassName}}{})
}
`

const createTableTemplate = `package migrations

import (
	"gorm.io/gorm"
)

// {{.ClassName}} migration - Create {{.TableName}} table
type {{.ClassName}} struct{}

// Up creates the {{.TableName}} table
func (m *{{.ClassName}}) Up(db *gorm.DB) error {
	return db.Exec(` + "`" + `
		CREATE TABLE {{.TableName}} (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			{{- range .Fields}}
			{{.Name}} {{toSQLType .Type}},
			{{- end}}
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	` + "`" + `).Error
}

// Down drops the {{.TableName}} table
func (m *{{.ClassName}}) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS {{.TableName}}").Error
}

// Description returns migration description
func (m *{{.ClassName}}) Description() string {
	return "Create {{.TableName}} table"
}

// Version returns migration version
func (m *{{.ClassName}}) Version() string {
	return "{{.Version}}"
}

// Auto-register migration
func init() {
	Register(&{{.ClassName}}{})
}
`

const alterTableTemplate = `package migrations

import (
	"gorm.io/gorm"
)

// {{.ClassName}} migration - Modify {{.TableName}} table
type {{.ClassName}} struct{}

// Up modifies the {{.TableName}} table
func (m *{{.ClassName}}) Up(db *gorm.DB) error {
	{{- range .Fields}}
	// Add {{.Name}} column
	if err := db.Exec("ALTER TABLE {{$.TableName}} ADD COLUMN {{.Name}} {{toSQLType .Type}}").Error; err != nil {
		return err
	}
	{{- end}}
	
	return nil
}

// Down reverts changes to the {{.TableName}} table
func (m *{{.ClassName}}) Down(db *gorm.DB) error {
	{{- range .Fields}}
	// Drop {{.Name}} column
	if err := db.Exec("ALTER TABLE {{$.TableName}} DROP COLUMN IF EXISTS {{.Name}}").Error; err != nil {
		return err
	}
	{{- end}}
	
	return nil
}

// Description returns migration description
func (m *{{.ClassName}}) Description() string {
	return "{{.Description}}"
}

// Version returns migration version
func (m *{{.ClassName}}) Version() string {
	return "{{.Version}}"
}

// Auto-register migration
func init() {
	Register(&{{.ClassName}}{})
}
`

const seederTemplate = `package seeders

import (
	"gorm.io/gorm"
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/logger"
	"go.uber.org/zap"
)

// {{.ClassName}} seeds the {{.TableName}} table
type {{.ClassName}} struct{}

// Run executes the seeder
func (s *{{.ClassName}}) Run(db *gorm.DB) error {
	logger.Info("Running {{.ClassName}}...")

	// TODO: Implement your seeding logic here
	// Example:
	// data := []entity.Model{
	//     {Field1: "value1", Field2: "value2"},
	//     {Field1: "value3", Field2: "value4"},
	// }
	//
	// return db.Create(&data).Error

	logger.Info("{{.ClassName}} completed successfully")
	return nil
}

// Name returns seeder name
func (s *{{.ClassName}}) Name() string {
	return "{{.ClassName}}"
}

// Auto-register seeder
func init() {
	Register(&{{.ClassName}}{})
}
`

const entityTemplate = `package entity

import (
	"time"

	"github.com/google/uuid"
	{{- if hasDecimalField .Fields}}
	"github.com/shopspring/decimal"
	{{- end}}
	"gorm.io/gorm"
)

// {{.EntityName}} represents a {{.EntityName}} entity
type {{.EntityName}} struct {
	ID        uuid.UUID      ` + "`json:\"id\" gorm:\"type:uuid;primary_key;default:gen_random_uuid()\"`" + `
	{{- range .Fields}}
	{{toPascalCase .Name}} {{toGoType .Type}} ` + "`json:\"{{.Name}}\" gorm:\"{{getGormTag .Type}}\" validate:\"{{getValidationTag .Type}}\"`" + `
	{{- end}}
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`json:\"-\" gorm:\"index\"`" + `
}

// Create{{.EntityName}}Request represents a request to create a {{.EntityName}}
type Create{{.EntityName}}Request struct {
	{{- range .Fields}}
	{{toPascalCase .Name}} {{toGoType .Type}} ` + "`json:\"{{.Name}}\" validate:\"{{getValidationTag .Type}}\"`" + `
	{{- end}}
}

// Update{{.EntityName}}Request represents a request to update a {{.EntityName}}
type Update{{.EntityName}}Request struct {
	{{- range .Fields}}
	{{toPascalCase .Name}} *{{toGoType .Type}} ` + "`json:\"{{.Name}},omitempty\" validate:\"omitempty,{{getValidationTag .Type}}\"`" + `
	{{- end}}
}

// {{.EntityName}}Filter represents filters for {{.EntityName}} queries
type {{.EntityName}}Filter struct {
	{{- range .Fields}}
	{{- if eq .Type "string"}}
	{{toPascalCase .Name}} string ` + "`form:\"{{.Name}}\"`" + `
	{{- end}}
	{{- end}}
	Search string ` + "`form:\"search\"`" + `
	Page   int    ` + "`form:\"page\" validate:\"min=1\"`" + `
	Limit  int    ` + "`form:\"limit\" validate:\"min=1,max=100\"`" + `
}

// TableName returns the table name for GORM
func ({{.EntityName}}) TableName() string {
	return "{{.TableName}}"
}
`

// Package templates - Simple structure without CRUD
const handlerTemplate = `package {{.PackageName}}

import (
	"go-clean-gin/pkg/errors"
	"go-clean-gin/pkg/logger"
	"go-clean-gin/pkg/response"
	"go-clean-gin/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type {{.EntityName}}Handler struct {
	usecase {{.EntityName}}Usecase
}

func New{{.EntityName}}Handler(usecase {{.EntityName}}Usecase) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{
		usecase: usecase,
	}
}

// TODO: Add your handler methods here
// Example:
// func (h *{{.EntityName}}Handler) SomeMethod(c *gin.Context) {
//     // Implementation here
// }
`

const portTemplate = `package {{.PackageName}}

import (
	"context"
)

// {{.EntityName}}Usecase defines the business logic interface for {{.PackageName}}
type {{.EntityName}}Usecase interface {
	// TODO: Add your usecase methods here
	// Example:
	// SomeMethod(ctx context.Context) error
}

// {{.EntityName}}Repository defines the data access interface for {{.PackageName}}
type {{.EntityName}}Repository interface {
	// TODO: Add your repository methods here
	// Example:
	// SomeMethod(ctx context.Context) error
}
`

const repositoryTemplate = `package {{.PackageName}}

import (
	"context"

	"gorm.io/gorm"
)

type {{.PackageName}}Repository struct {
	db *gorm.DB
}

func New{{.EntityName}}Repository(db *gorm.DB) {{.EntityName}}Repository {
	return &{{.PackageName}}Repository{
		db: db,
	}
}

// TODO: Add your repository methods here
// Example:
// func (r *{{.PackageName}}Repository) SomeMethod(ctx context.Context) error {
//     return r.db.WithContext(ctx).Error
// }
`

const usecaseTemplate = `package {{.PackageName}}

import (
	"context"
	"go-clean-gin/pkg/errors"
	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
)

type {{.PackageName}}Usecase struct {
	repo {{.EntityName}}Repository
}

func New{{.EntityName}}Usecase(repo {{.EntityName}}Repository) {{.EntityName}}Usecase {
	return &{{.PackageName}}Usecase{
		repo: repo,
	}
}

// TODO: Add your usecase methods here
// Example:
// func (u *{{.PackageName}}Usecase) SomeMethod(ctx context.Context) error {
//     logger.Info("Executing SomeMethod for {{.PackageName}}")
//     
//     if err := u.repo.SomeMethod(ctx); err != nil {
//         logger.Error("Failed to execute SomeMethod", zap.Error(err))
//         return errors.Wrap(err, errors.ErrInternal, "Failed to execute SomeMethod", 500)
//     }
//     
//     return nil
// }
`
