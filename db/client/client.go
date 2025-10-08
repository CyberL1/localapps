package dbClient

import (
	"database/sql"
	"fmt"
	"localapps/constants"
	db "localapps/db/generated"
	"localapps/resources"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

func Migrate() {
	if err := goose.SetDialect("sqlite"); err != nil {
		fmt.Println(err)
	}

	sql, err := goose.OpenDBWithDriver("sqlite", filepath.Join(constants.LocalappsDir, "localapps.db"))
	if err != nil {
		fmt.Println(err)
	}

	goose.SetBaseFS(resources.Resources)

	if err := goose.Up(sql, "db_migrations"); err != nil {
		fmt.Println(err)
	}
}

func GetClient() (*db.Queries, error) {
	conn, err := sql.Open("sqlite", filepath.Join(constants.LocalappsDir, "localapps.db"))
	if err != nil {
		return nil, err
	}

	client := db.New(conn)
	return client, nil
}
