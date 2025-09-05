package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func connectToDb() error {
	var err error
	dbConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}
	return nil
}
