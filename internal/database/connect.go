package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

var dbInstance *DB

func ConnectDB() *DB {
	if dbInstance != nil {
		return dbInstance
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, db_name)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}

	dbInstance = &DB{Pool: pool}
	fmt.Println("Connected to database")
	return dbInstance
}

func GetDB() *DB {
	return dbInstance
}
