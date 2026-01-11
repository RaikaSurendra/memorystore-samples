package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	// Default values for local testing if env vars not set (BE CAREFUL IN PROD)
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbName == "" {
		dbName = "postgres"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPass, dbHost, dbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Table creation is handled by init.sql in Docker environment

	log.Println("Connected to PostgreSQL")
}

func GetItemFromDB(id int64) (string, string, float64, error) {
	var name, description string
	var price float64
	err := Pool.QueryRow(context.Background(), "SELECT name, description, price FROM items WHERE id=$1", id).Scan(&name, &description, &price)
	return name, description, price, err
}

func AddItemToDB(name string, description string, price float64) (int64, error) {
	var id int64
	err := Pool.QueryRow(context.Background(), "INSERT INTO items (name, description, price) VALUES ($1, $2, $3) RETURNING id", name, description, price).Scan(&id)
	return id, err
}

func DeleteItemFromDB(id int64) error {
	_, err := Pool.Exec(context.Background(), "DELETE FROM items WHERE id=$1", id)
	return err
}
