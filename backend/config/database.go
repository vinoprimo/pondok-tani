package config

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"context"
)

var DB *pgx.Conn

func ConnectDB() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	DB = conn
	fmt.Println("Database connected!")
}