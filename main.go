package main

import (
	"log"
	"net/http"

	"github.com/SWernerH/test_1/db"
	"github.com/SWernerH/test_1/routes"
)

func main() {
	db.InitDB()
	routes.RegisterRoutes()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}