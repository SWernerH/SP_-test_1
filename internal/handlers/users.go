package handlers

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/SWernerH/test_1/internal/db"
    "github.com/SWernerH/test_1/internal/models"
    "github.com/SWernerH/test_1/internal/validator"
)

type Application struct {
    DB *db.DB
}

// JSON helpers
type envelope map[string]any

func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope) {
    js, err := json.Marshal(data)
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
    http.Error(w, "Server error: "+err.Error(), 500)
}

func (app *Application) notFound(w http.ResponseWriter) {
    http.Error(w, "Resource not found", 404)
}

// ─── Handlers ───────────────────────────────────────────────

// GET /users
func (app *Application) ListUsers(w http.ResponseWriter, r *http.Request) {
    query := `SELECT id, username, email FROM users ORDER BY id`
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    rows, err := app.DB.QueryContext(ctx, query)
    if err != nil {
        app.serverError(w, err)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
            app.serverError(w, err)
            return
        }
        users = append(users, u)
    }

    app.writeJSON(w, http.StatusOK, envelope{"users": users})
}

// GET /users/{id}
func (app *Application) GetUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    query := `SELECT id, username, email FROM users WHERE id = $1`
    var u models.User

    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    err = app.DB.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    app.writeJSON(w, http.StatusOK, envelope{"user": u})
}

// POST /users
func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
    var input models.User
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", 400)
        return
    }

    v := validator.New()
    v.Check(input.Username != "", "username", "must be provided")
    v.Check(input.Email != "", "email", "must be provided")
    v.Check(input.Password != "", "password", "must be provided")

    if !v.Valid() {
        app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        app.serverError(w, err)
        return
    }

    query := `INSERT INTO users (username, email, password_hash, created_at)
              VALUES ($1, $2, $3, NOW()) RETURNING id`

    var newID int64
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    err = app.DB.QueryRowContext(ctx, query, input.Username, input.Email, string(hash)).Scan(&newID)
    if err != nil {
        app.serverError(w, err)
        return
    }

    input.ID = newID
    app.writeJSON(w, http.StatusCreated, envelope{"user": input})
}

// PUT /users/{id}
func (app *Application) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    var input models.User
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", 400)
        return
    }

    query := `UPDATE users SET username=$1, email=$2 WHERE id=$3`
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    result, err := app.DB.ExecContext(ctx, query, input.Username, input.Email, id)
    if err != nil {
        app.serverError(w, err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        app.serverError(w, err)
        return
    }
    if rowsAffected == 0 {
        app.notFound(w)
        return
    }

    input.ID = id
    app.writeJSON(w, http.StatusOK, envelope{"user": input})
}

// DELETE /users/{id}
func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    query := `DELETE FROM users WHERE id=$1`
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    result, err := app.DB.ExecContext(ctx, query, id)
    if err != nil {
        app.serverError(w, err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        app.serverError(w, err)
        return
    }
    if rowsAffected == 0 {
        app.notFound(w)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// GET /health
func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    err := app.DB.PingContext(ctx)
    dbStatus := "reachable"
    if err != nil {
        dbStatus = "unreachable: " + err.Error()
    }

    app.writeJSON(w, http.StatusOK, envelope{
        "status":   "available",
        "database": dbStatus,
    })
}
