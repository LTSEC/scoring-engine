package scoring

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Tables map[string]map[string]interface{} `yaml:"tables"`
}

func createConStr(dbType string, address string, portNum int, username string, password string, DBName string) (string, error) {
	switch dbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, address, portNum, DBName), nil
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			address, portNum, username, password, DBName), nil
	default:
		return "", fmt.Errorf("invalid DB parameter: %s", dbType)
	}
}

func createConStr_sqlite(DBName string) string {
	return filepath.Join("database_files", DBName)
}

func getTableNames(db *sql.DB, dbType string) ([]string, error) {
	var query string
	switch dbType {
	case "sqlite3":
		query = `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`
	case "postgres":
		query = `SELECT tablename FROM pg_tables WHERE schemaname='public'`
	case "mysql":
		query = `SHOW TABLES`
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("error scanning table name: %v", err)
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func normalizeQuery(query string, dbType string) string {
	switch dbType {
	case "postgres":
		return strings.ReplaceAll(query, "\"", "")
	case "mysql":
		return strings.ReplaceAll(query, "`", "")
	default:
		return query
	}
}

func compareWithSQLite(sourceDB *sql.DB, sqliteDB *sql.DB, tableName string, sourceType string) error {
	// Get SQLite schema and data first as reference
	sqliteSchema, err := sqliteDB.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return fmt.Errorf("error getting SQLite schema: %v", err)
	}
	defer sqliteSchema.Close()

	var sqliteColumns []string
	var sqliteTypes []string
	for sqliteSchema.Next() {
		var cid, notnull, pk int
		var name, dtype, dflt_value sql.NullString
		if err := sqliteSchema.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("error scanning SQLite schema: %v", err)
		}
		sqliteColumns = append(sqliteColumns, name.String)
		sqliteTypes = append(sqliteTypes, dtype.String)
	}

	// Prepare and execute queries
	sqliteQuery := fmt.Sprintf("SELECT * FROM %s", tableName)
	sourceQuery := normalizeQuery(sqliteQuery, sourceType)

	sqliteRows, err := sqliteDB.Query(sqliteQuery)
	if err != nil {
		return fmt.Errorf("error querying SQLite: %v", err)
	}
	defer sqliteRows.Close()

	sourceRows, err := sourceDB.Query(sourceQuery)
	if err != nil {
		return fmt.Errorf("error querying source database: %v", err)
	}
	defer sourceRows.Close()

	// Compare data
	sqliteData := make([][]interface{}, 0)
	sourceData := make([][]interface{}, 0)

	// Read SQLite data
	for sqliteRows.Next() {
		values := make([]interface{}, len(sqliteColumns))
		valuePtrs := make([]interface{}, len(sqliteColumns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		if err := sqliteRows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("error scanning SQLite row: %v", err)
		}

		// Normalize SQLite types
		for i, v := range values {
			if v == nil {
				continue
			}
			switch sqliteTypes[i] {
			case "INTEGER":
				if n, ok := v.(int64); ok {
					values[i] = n
				}
			case "REAL":
				if f, ok := v.(float64); ok {
					values[i] = f
				}
			case "TEXT":
				if b, ok := v.([]byte); ok {
					values[i] = string(b)
				}
			case "DATETIME":
				if t, ok := v.(time.Time); ok {
					values[i] = t.Format("2006-01-02 15:04:05")
				} else if s, ok := v.(string); ok {
					values[i] = s
				}
			}
		}
		sqliteData = append(sqliteData, values)
	}

	// Read source data with SQLite types as reference
	sourceCols, _ := sourceRows.Columns()
	for sourceRows.Next() {
		values := make([]interface{}, len(sourceCols))
		valuePtrs := make([]interface{}, len(sourceCols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		if err := sourceRows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("error scanning source row: %v", err)
		}

		// Normalize source types to match SQLite
		for i, v := range values {
			if v == nil {
				continue
			}
			switch sqliteTypes[i] {
			case "INTEGER":
				switch val := v.(type) {
				case []byte:
					values[i], _ = parseInt(string(val))
				case int64:
					values[i] = val
				case int32:
					values[i] = int64(val)
				}
			case "REAL":
				switch val := v.(type) {
				case []byte:
					values[i], _ = parseFloat(string(val))
				case float64:
					values[i] = val
				case float32:
					values[i] = float64(val)
				}
			case "TEXT":
				switch val := v.(type) {
				case []byte:
					values[i] = string(val)
				case string:
					values[i] = val
				}
			case "DATETIME":
				switch val := v.(type) {
				case []byte:
					values[i] = string(val)
				case time.Time:
					values[i] = val.Format("2006-01-02 15:04:05")
				case string:
					values[i] = val
				}
			}
		}
		sourceData = append(sourceData, values)
	}

	// Compare row counts
	if len(sqliteData) != len(sourceData) {
		return fmt.Errorf("row count mismatch: sqlite=%d, source=%d", len(sqliteData), len(sourceData))
	}

	// Compare data
	for i := range sqliteData {
		for j := range sqliteData[i] {
			if fmt.Sprintf("%v", sqliteData[i][j]) != fmt.Sprintf("%v", sourceData[i][j]) {
				return fmt.Errorf("data mismatch at row %d, column %s: sqlite=%v, source=%v",
					i+1, sqliteColumns[j], sqliteData[i][j], sourceData[i][j])
			}
		}
	}

	return nil
}

func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func DBcheck(dbType string, address string, portNum int, username string, password string, DBName string, DBFile string) (bool, error) {
	// Connect to SQLite database first (reference)
	sqliteDB, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		return false, fmt.Errorf("failed to open SQLite database: %v", err)
	}
	defer sqliteDB.Close()

	// Connect to source database
	sourceConnStr, err := createConStr(dbType, address, portNum, username, password, DBName)
	if err != nil {
		return false, fmt.Errorf("failed to create source connection string: %v", err)
	}

	sourceDB, err := sql.Open(dbType, sourceConnStr)
	if err != nil {
		return false, fmt.Errorf("failed to open source database: %v", err)
	}
	defer sourceDB.Close()

	// Test source connection
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	if err := sourceDB.PingContext(ctx); err != nil {
		return false, fmt.Errorf("source connection test failed: %v", err)
	}

	// Get table names from SQLite
	tables, err := getTableNames(sqliteDB, "sqlite3")
	if err != nil {
		return false, fmt.Errorf("failed to get SQLite table names: %v", err)
	}

	// Compare each table
	for _, tableName := range tables {
		if err := compareWithSQLite(sourceDB, sqliteDB, tableName, dbType); err != nil {
			return false, fmt.Errorf("table %s comparison failed: %v", tableName, err)
		}
	}

	return true, nil
}

func ScoreDB(dbType string, address string, portNum int, username string, password string, DBName string, DBFile string) (int, bool, error) {
	success, err := DBcheck(dbType, address, portNum, username, password, DBName, DBFile)
	if err != nil {
		return 0, false, fmt.Errorf("DB scoring failed: %v", err)
	}
	return 1, success, nil
}