package main

import (
	"database/sql"
	"fmt"
	"log"
	"serverAPI/cmd/api"
	"serverAPI/config"
	"serverAPI/repo"
)

func main() {
	db, _ := repo.NewPostgresDB(
		fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			config.Env.DBUser,
			config.Env.DBPassword,
			config.Env.Host,
			config.Env.Port,
			config.Env.DBName,
		),
	)
	initStorage(db)
	server := api.NewAPIServer("localhost:8000", db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run a server")
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB is successfully connected")
}
