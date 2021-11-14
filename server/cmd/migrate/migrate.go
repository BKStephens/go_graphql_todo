package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
)

func main() {
	var databaseUrl string
	if envVar := os.Getenv("DATABASE_URL"); envVar != "" {
		databaseUrl = envVar
	} else {
		panic("DATABASE_URL not set")
	}
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Println("Couldn't connect to the database")
		fmt.Println(err)
		return
	}
	var m *migrate.Migrator
	m, err = migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		fmt.Printf("Unable to create migrator: %v", err)
		return
	}

	err = m.LoadMigrations("server/db/migrations")
	if err != nil {
		fmt.Println("Cannot load migrations. Make sure you are running this command from the root of the repository.")
		fmt.Println(err)
	}

	err = m.Migrate(context.Background())
	if err != nil {
		fmt.Println("Migration failed")
		fmt.Println(err)
		return
	}

	m.OnStart = func(_ int32, name, direction, _ string) {
		fmt.Printf("Migrating %s: %s", direction, name)
	}
}
