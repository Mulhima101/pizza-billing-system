package config

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

type Database struct {
    DB *sql.DB
}

func NewDatabase() (*Database, error) {
    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "5432")
    user := getEnv("DB_USER", "postgres")
    password := getEnv("DB_PASSWORD", "your_password")
    dbname := getEnv("DB_NAME", "pizza_billing")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }

    return &Database{DB: db}, nil
}

func (d *Database) Close() error {
    return d.DB.Close()
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}