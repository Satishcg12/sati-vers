package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/satishcg12/sati-vers/auth/user-service/config"
)

type Database struct {
	cfg config.Database
}

func NewDatabase(cfg config.Database) *Database {
	return &Database{
		cfg: cfg,
	}
}

func (d *Database) Connect() (*sql.DB, error) {
	pgsqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.cfg.Host, d.cfg.Port, d.cfg.User, d.cfg.Password, d.cfg.Name)

	db, err := sql.Open("postgres", pgsqlconn)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}

func (d *Database) Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v\n", err)
	}
	log.Println("Closed database connection")
}
