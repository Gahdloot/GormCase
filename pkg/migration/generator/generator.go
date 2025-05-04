package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ModelInfo holds information about a model
type ModelInfo struct {
	Name       string
	Fields     []FieldInfo
	TableName  string
	ImportPath string
}

// FieldInfo holds information about a model field
type FieldInfo struct {
	Name       string
	Type       string
	IsNullable bool
	IsUnique   bool
	IsPrimary  bool
	Default    string
}

// GenerateMigration creates a migration file based on model changes
func GenerateMigration(modelPath string) error {
	// Parse the model file
	modelInfo, err := parseModelFile(modelPath)
	if err != nil {
		return fmt.Errorf("failed to parse model file: %v", err)
	}

	// Generate migration content
	migrationContent := generateMigrationContent(modelInfo)

	// Create migrations directory if it doesn't exist
	migrationsDir := filepath.Join(filepath.Dir(modelPath), "..", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	// Generate migration filename
	timestamp := time.Now().Format("20060102150405")
	migrationName := fmt.Sprintf("%s_create_%s_table", timestamp, strings.ToLower(modelInfo.TableName))
	migrationPath := filepath.Join(migrationsDir, migrationName+".go")

	// Write migration file
	if err := os.WriteFile(migrationPath, []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %v", err)
	}

	fmt.Printf("Generated migration file: %s\n", migrationPath)
	return nil
}

func parseModelFile(path string) (*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelInfo := &ModelInfo{
		Fields: make([]FieldInfo, 0),
	}

	// Find the model struct
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			modelInfo.Name = typeSpec.Name.Name

			// Parse fields
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				fieldInfo := FieldInfo{
					Name: field.Names[0].Name,
				}

				// Parse field type
				switch t := field.Type.(type) {
				case *ast.StarExpr:
					if sel, ok := t.X.(*ast.SelectorExpr); ok {
						fieldInfo.Type = sel.Sel.Name
						fieldInfo.IsNullable = true
					}
				case *ast.Ident:
					fieldInfo.Type = t.Name
				}

				// Parse tags
				if field.Tag != nil {
					tag := field.Tag.Value
					fieldInfo.IsUnique = strings.Contains(tag, "unique")
					fieldInfo.IsPrimary = strings.Contains(tag, "primary")
				}

				modelInfo.Fields = append(modelInfo.Fields, fieldInfo)
			}
		}
	}

	// Find TableName method
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "TableName" {
			continue
		}

		if stmt, ok := funcDecl.Body.List[0].(*ast.ReturnStmt); ok {
			if lit, ok := stmt.Results[0].(*ast.BasicLit); ok {
				modelInfo.TableName = strings.Trim(lit.Value, "\"")
			}
		}
	}

	return modelInfo, nil
}

func generateMigrationContent(model *ModelInfo) string {
	var fields []string
	for _, field := range model.Fields {
		fieldDef := fmt.Sprintf("%s %s", strings.ToLower(field.Name), getSQLType(field.Type))
		if !field.IsNullable {
			fieldDef += " NOT NULL"
		}
		if field.IsUnique {
			fieldDef += " UNIQUE"
		}
		if field.IsPrimary {
			fieldDef += " PRIMARY KEY"
		}
		fields = append(fields, fieldDef)
	}

	tableName := strings.ToLower(model.TableName)
	if tableName == "" {
		tableName = strings.ToLower(model.Name) + "s"
	}

	template := `package migrations

import (
	"context"
	"database/sql"
	"fmt"
)

// Migration%s represents the migration for %s
type Migration%s struct{}

// Up creates the %s table
func (m *Migration%s) Up(ctx context.Context, tx *sql.Tx) error {
	query := "CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, %s, created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	_, err := tx.ExecContext(ctx, fmt.Sprintf(query, "%s", "%s"))
	return err
}

// Down drops the %s table
func (m *Migration%s) Down(ctx context.Context, tx *sql.Tx) error {
	query := "DROP TABLE IF EXISTS %s CASCADE;"
	_, err := tx.ExecContext(ctx, fmt.Sprintf(query, "%s"))
	return err
}

// Name returns the name of the migration
func (m *Migration%s) Name() string {
	return "%s"
}`

	return fmt.Sprintf(template,
		strings.Title(tableName),
		model.Name,
		strings.Title(tableName),
		tableName,
		strings.Title(tableName),
		tableName,
		strings.Join(fields, ", "),
		tableName,
		strings.Title(tableName),
		tableName,
		strings.Title(tableName),
		time.Now().Format("20060102150405")+"_create_"+tableName+"_table",
	)
}

func getSQLType(goType string) string {
	switch goType {
	case "CharField":
		return "VARCHAR(255)"
	case "TextField":
		return "TEXT"
	case "IntegerField":
		return "INTEGER"
	case "BooleanField":
		return "BOOLEAN"
	case "DateTimeField":
		return "TIMESTAMP WITH TIME ZONE"
	default:
		return "VARCHAR(255)"
	}
}
