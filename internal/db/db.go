package db

import (
    "context"
    "database/sql"
    "log"
    "time"

    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

func OpenDB(dsn string) (*DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxIdleTime(15 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, err
    }

    log.Println("Database connection pool established")
    return &DB{db}, nil
}
