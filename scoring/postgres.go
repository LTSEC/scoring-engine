package scoring

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// PGconnect connects to a PostgreSQL database on the target host and verifies
// the schema by running a random query from the provided query file.
func PGconnect(address string, portNum int, username string, password string, DBName string, DBPath string) (bool, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=1",
		address, portNum, username, password, DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return false, err
	}

	// Verify the database by running a random query from the query file
	isValid, err := PGverify(db, DBPath)
	if err != nil {
		return false, fmt.Errorf("database verification failed: %w", err)
	}

	if !isValid {
		return false, fmt.Errorf("database verification returned invalid")
	}

	return true, nil
}

// PGverify reads a query file, picks a random SQL statement, and executes it
// against the already-open PostgreSQL connection.
func PGverify(db *sql.DB, DBPath string) (bool, error) {
	content, err := os.ReadFile(DBPath)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var validLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			validLines = append(validLines, trimmed)
		}
	}

	if len(validLines) == 0 {
		return false, fmt.Errorf("no valid SQL statements found in the file")
	}

	query := validLines[rand.Intn(len(validLines))]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}

	return true, nil
}

// ScorePostgres checks a PostgreSQL service on the target host and assigns points.
func ScorePostgres(address string, portNum int, username string, password string, DBName string, DBPath string) (int, bool, error) {
	_, err := PGconnect(address, portNum, username, password, DBName, DBPath)
	if err != nil {
		return 0, false, fmt.Errorf("PostgreSQL scoring failed: %v", err)
	}
	return successPoints, true, nil
}
