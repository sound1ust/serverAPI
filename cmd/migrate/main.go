package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"serverAPI/config"
)

type Logger struct{}

func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Logger) Verbose() bool {
	return true
}

func main() {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.Host,
		config.Env.Port,
		config.Env.DBName,
	)
	m, err := migrate.New("file://cmd/migrate/migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}
	m.Log = &Logger{}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		fmt.Println("Successfully applied migrations!")
	} else if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		fmt.Printf("Successfully reverted migrations!")
	} else {
		log.Fatalf("Unknown action: %s", cmd)
	}
}
