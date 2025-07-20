package config

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

type Database struct {
    DB *sql.DB
}

func NewDatabase() (*Database, error) {
    connStr := "user=postgres dbname=pizza_billing sslmode=disable password=your_password"
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