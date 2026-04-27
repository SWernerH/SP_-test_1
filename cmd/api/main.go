package main

import (
    "log"
    "net/http"
    "os"

    "streaming_platform/internal/db"
    "streaming_platform/internal/handlers"
)

func main() {
    dsn := os.Getenv("STREAMING_DSN")
    if dsn == "" {
        log.Fatal("STREAMING_DSN environment variable not set")
    }

    conn, err := db.OpenDB(dsn)
    if err != nil {
        log.Fatalf("cannot open database: %v", err)
    }
    defer conn.Close()

    app := &handlers.Application{DB: conn}

    mux := http.NewServeMux()
    mux.HandleFunc("GET /users", app.ListUsers)
    mux.HandleFunc("GET /users/{id}", app.GetUser)
    mux.HandleFunc("POST /users", app.CreateUser)
    mux.HandleFunc("PUT /users/{id}", app.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", app.DeleteUser)
    mux.HandleFunc("GET /health", app.Health)

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
