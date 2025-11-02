package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"auth/config"

	_ "github.com/lib/pq"
)

// var DB *sql.DB

// func Connect() {
// 	connStr := config.GetEnv("DATABASE_URL", "")
// 	if connStr == "" {
// 		log.Fatal("DATABASE_URL not found in environment")
// 	}

// 	var err error
// 	DB, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to open database:", err)
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	DB.SetMaxOpenConns(10) // Vercel allows ~10
//     DB.SetMaxIdleConns(2)

// 	fmt.Println("Connected to PostgreSQL successfully!")
// }

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		connStr := config.GetEnv("DATABASE_URL", "")
		if connStr == "" {
			log.Fatal("DATABASE_URL not found in environment")
		}

		var err error
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}

		err = DB.Ping()
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		DB.SetMaxOpenConns(10)
		DB.SetMaxIdleConns(2)
		DB.SetConnMaxLifetime(0) // Let Vercel manage

		fmt.Println("Connected to PostgreSQL successfully!")
	})
}

// Helper to get DB (use in handlers)
func GetDB() *sql.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}