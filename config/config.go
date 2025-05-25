package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Custom error types for better error handling
type ConfigError struct {
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	return e.Message
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// sanitizeError removes sensitive information from error messages
func sanitizeError(err error) error {
	if err == nil {
		return nil
	}

	// Convert error to string and check for sensitive patterns
	errStr := err.Error()

	// Remove any potential connection strings
	if strings.Contains(errStr, "@tcp") {
		return fmt.Errorf("database connection error")
	}

	// Remove any potential credentials
	if strings.Contains(errStr, "Access denied") {
		return fmt.Errorf("database authentication error")
	}

	// Remove any potential host information
	if strings.Contains(errStr, "no such host") {
		return fmt.Errorf("database host error")
	}

	// For other errors, return a generic message
	return fmt.Errorf("database error")
}

func InitDB() (*sql.DB, error) {
	// Get database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, &ConfigError{Message: "database user not configured"}
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		return nil, &ConfigError{Message: "database password not configured"}
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, &ConfigError{Message: "database host not configured"}
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, &ConfigError{Message: "database port not configured"}
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, &ConfigError{Message: "database name not configured"}
	}

	// Create database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, &ConfigError{
			Message: "failed to initialize database connection",
			Err:     sanitizeError(err),
		}
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, &ConfigError{
			Message: "failed to connect to database",
			Err:     sanitizeError(err),
		}
	}

	return db, nil
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
