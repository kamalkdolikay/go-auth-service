// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"auth/config"
	"auth/db"
	"auth/handlers"
	"auth/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// 1. Load .env file (only in local dev)
	config.LoadEnv()

	// 2. Initialize JWT secret from env
	handlers.InitJWT()

	// 3. Initialize DB connection
	db.InitDB()

	// 4. Initialize DB schema (idempotent)
	if err := initDBSchema(); err != nil {
		log.Fatal("Failed to initialize DB schema:", err)
	}

	// 5. Setup routes
	router := mux.NewRouter()
	routes.RegisterRoutesToMux(router)

	// 6. CORS middleware (wrap router)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins for dev
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, 
		AllowCredentials: true,
	})

	handler := c.Handler(router) // wrap the router

	// 7. Start server
	port := config.GetEnv("PORT", "8000")
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// initDBSchema creates the users table if it doesn't exist
func initDBSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	_, err := db.GetDB().Exec(query)
	return err
}