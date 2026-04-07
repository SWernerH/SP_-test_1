package routes

import (
	"net/http"

	"github.com/SWernerH/test_1/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUsers(w, r)
		case http.MethodPost:
			handlers.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})
}