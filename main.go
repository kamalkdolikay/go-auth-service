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
	// 1. Load .env file
	config.LoadEnv()

	// 2. Initialize JWT secret
	handlers.InitJWT()

	// 3. Initialize DB connection
	db.InitDB()

	// 4. Initialize DB schema (Updated for Admin & Activity)
	if err := initDBSchema(); err != nil {
		log.Fatal("Failed to initialize DB schema:", err)
	}

	// 5. Setup routes
	router := mux.NewRouter()
	routes.RegisterRoutesToMux(router)

	// 6. CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// 7. Start server
	port := config.GetEnv("PORT", "8000")
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// initDBSchema creates tables for Users and UserActivity
func initDBSchema() error {
	dbConn := db.GetDB()

	// 1. Users Table
	userQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role VARCHAR(50) DEFAULT 'user',
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	if _, err := dbConn.Exec(userQuery); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// 2. Activity Table (For metrics)
	activityQuery := `
		CREATE TABLE IF NOT EXISTS user_activity (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			prompt TEXT,
			topic_detected TEXT,
			target_language VARCHAR(10),
			request_type VARCHAR(10),
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	if _, err := dbConn.Exec(activityQuery); err != nil {
		return fmt.Errorf("error creating activity table: %w", err)
	}

	return nil
}
