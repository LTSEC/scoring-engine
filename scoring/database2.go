package scoring

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// for the merge, this code may not be needed. this is mainly for db_test.go.

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	DBType   string
	DBRef    string
}

// Creates a database
func CreateDatabase(cfg Config) error {
	err := ManageDatabase(cfg, "create")
	return err
}

// Drops a database
func DropDatabase(cfg Config) error {
	err := ManageDatabase(cfg, "drop")
	return err
}

// General method for creating/dropping a database (to reuse code)
func ManageDatabase(cfg Config, action string) error {
	connStr, err := createConStr(cfg.DBType, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBType)
	if err != nil {
		return fmt.Errorf("failed to create connection string for db type %s: %w", cfg.DBType, err)
	}

	db, err := sql.Open(cfg.DBType, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", cfg.DBName, err)
	}
	defer db.Close()

	// Check if database exists using DB specific syntax
	var exists bool
	if cfg.DBType == "postgres" {
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DBName).Scan(&exists)
	} else {
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)", cfg.DBName).Scan(&exists)
	}
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	var query string

	switch action {
	// Invert the exists flag because we want to do the create a database if it doesn't exist
	case "create":
		query = fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		exists = !exists
	// No need to invert it because we want to drop a database if it exists
	case "drop":
		query = fmt.Sprintf("DROP DATABASE %s", cfg.DBName)
	// Invalid action returns an error
	default:
		return fmt.Errorf("Invalid Action")
	}

	// Perform a create/drop depending on the exists flag
	if exists {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}
	return nil
}

// For setting up the schema of a DB and seeding it
func ExecuteSQLFile(cfg Config, fileName string) error {
	// Assuming that the test is in scoring, the files are stored in database_files which is a directory in scoring
	filePath := filepath.Join("database_files", fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Errorf("File doesn't exist")
	}

	// Read SQL file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Connect to the specific database
	connStr, err := createConStr(cfg.DBType, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	if err != nil {
		return fmt.Errorf("failed to create connection string: %w", err)
	}

	// Open a connection to the DB given the driver and earlier connection string
	db, err := sql.Open(cfg.DBType, connStr)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer db.Close()

	// Split the content into individual statements
	statements := splitSQLStatements(string(content))

	// Execute each statement
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		// Handle PostgreSQL dollar-quoted string literals
		if cfg.DBType == "postgres" {
			stmt = strings.ReplaceAll(stmt, "$$", "'")
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement: %w\nStatement: %s", err, stmt)
		}

	}

	return nil
}

// Sanitizes the SQL file before it gets run. This is so we can use the same function for creating a schema and seeding
func splitSQLStatements(content string) []string {
	var statements []string
	var currentStmt strings.Builder

	// State tracking
	inString := false
	inLineComment := false
	inBlockComment := false

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		lineLen := len(line)
		skipChar := false

		for j := 0; j < lineLen; j++ {
			if skipChar {
				skipChar = false
				continue
			}

			char := line[j]

			// Handle string literals
			if char == '\'' && !inLineComment && !inBlockComment {
				inString = !inString
			}

			// Handle line comments --
			if !inString && !inBlockComment && j < lineLen-1 && char == '-' && line[j+1] == '-' {
				inLineComment = true
				skipChar = true
				continue
			}

			// Handle block comments /* */
			if !inString && !inLineComment && j < lineLen-1 && char == '/' && line[j+1] == '*' {
				inBlockComment = true
				skipChar = true
				continue
			}
			if !inString && inBlockComment && j < lineLen-1 && char == '*' && line[j+1] == '/' {
				inBlockComment = false
				skipChar = true
				continue
			}

			// Only add characters if we're not in a comment
			if !inLineComment && !inBlockComment {
				// Add the character to current statement
				currentStmt.WriteByte(char)

				// Check for statement end
				if char == ';' && !inString {
					stmt := strings.TrimSpace(currentStmt.String())
					if stmt != "" {
						statements = append(statements, stmt)
					}
					currentStmt.Reset()
				}
			}
		}

		// Reset line comment flag at end of line
		inLineComment = false

		// Add newline if not at the end of all lines and not in a comment
		if i < len(lines)-1 && !inBlockComment {
			currentStmt.WriteString("\n")
		}
	}

	// Add final statement if it doesn't end with semicolon
	finalStmt := strings.TrimSpace(currentStmt.String())
	if finalStmt != "" {
		statements = append(statements, finalStmt)
	}

	// Clean up statements by removing empty ones and extra whitespace
	var cleanStatements []string
	for _, stmt := range statements {
		cleaned := strings.TrimSpace(stmt)
		if cleaned != "" {
			cleanStatements = append(cleanStatements, cleaned)
		}
	}

	return cleanStatements
}
