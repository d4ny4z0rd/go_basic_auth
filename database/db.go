package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err!=nil {
		log.Fatalf("Error loading the env file: %v \n", err)
	}

	port, host, user, name, password := os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")
	fmt.Println("Loaded environment variables for database")
	
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, name, user, password)
	db, err := sql.Open("postgres", connectionString)
	if err!=nil {
		log.Fatalf("Error connecting to the database: %v \n", err)
	}

	if err = db.Ping(); err!=nil {
		log.Fatalf("Database ping failed: %v \n", err)
	}

	println("Connected to the database")

	return db
}

func CreateUsersTable(db *sql.DB) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			firstname	VARCHAR(100)	NOT NULL,
			lastname	VARCHAR(100)	NOT NULL,
			password	VARCHAR(100)	NOT NULL,
			email		VARCHAR(100)	NOT NULL	UNIQUE,
			userid		UUID			NOT NULL	UNIQUE
		);
	`

	_, err := db.Exec(createTableQuery)
	if err!=nil {
		log.Fatalf("Error creating users table: %v \n", err)
	}

	log.Println("Users table created successfully")
}
