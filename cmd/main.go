package main

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/FarStep131/go-jwt/docker/ent"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres dbname=mydb password=password sslmode=disable"
	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		panic(err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		panic(err)
	}

	fmt.Print("migrated")
}
