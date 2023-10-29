package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"github.com/FarStep131/go-jwt/docker/ent"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBClient struct {
	Client *ent.Client
}

func NewDBClient() (*DBClient, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	host := os.Getenv("DBHost")
	port := os.Getenv("DBPort")
	user := os.Getenv("DBUser")
	database := os.Getenv("DBName")
	password := os.Getenv("DBPassword")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, database, password)
	db, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Schema.Create(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return &DBClient{Client: db}, nil
}
