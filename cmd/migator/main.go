package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	var databaseURL string
	var migrationsPath string
	flag.StringVar(&databaseURL, "database-url", "", "database url")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if databaseURL == "" {
		panic("database url is required")
	}

	if migrationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New("file://"+migrationsPath, databaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("applied migrations successfully")
}
