package dbClient

import (
	"context"
	"fmt"
	db "localapps/db/generated"
	"localapps/resources"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
)

var client *db.Queries

func Migrate() {
	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Println(err)
	}

	sql, err := goose.OpenDBWithDriver("pgx", os.Getenv("LOCALAPPS_DB"))
	if err != nil {
		fmt.Println(err)
	}

	goose.SetBaseFS(resources.Resources)

	if err := goose.Up(sql, "db_migrations"); err != nil {
		fmt.Println(err)
	}
}

func GetClient() (*db.Queries, error) {
	if client == nil {
		conn, err := pgx.Connect(context.Background(), os.Getenv("LOCALAPPS_DB"))
		if err != nil {
			return nil, err
		}

		client = db.New(conn)
	}
	return client, nil
}
