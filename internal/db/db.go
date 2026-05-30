package db

import (
	"database/sql"
	"embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	if err := goose.Up(DB, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func createTables() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);`

	todosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN NOT NULL DEFAULT 0
	);`

	if _, err := DB.Exec(usersTable); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	if _, err := DB.Exec(todosTable); err != nil {
		log.Fatalf("Failed to create todos table: %v", err)
	}
}
