package api

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}


func ConnectToDB() (*sql.DB, error) {
	var cfg DBConfig
	cfg.Host = "localhost"
	cfg.Port = 5432
	cfg.User = "postgres"
	cfg.Password = "1234"
	cfg.DBName = "pendek.in"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Try to connect
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}
