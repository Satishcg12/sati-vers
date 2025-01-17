package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/satishcg12/sati-vers/auth/authorization-service/config"
)

type Database struct {
	dbCfg config.Database
}

type IDatabase interface {
	Connect() (*sql.DB, error)
	AutoMigrate(*sql.DB, string) error
}

func NewDatabase(cfg config.Database) *Database {
	return &Database{
		dbCfg: cfg,
	}
}

func (d *Database) Connect() (*sql.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.dbCfg.Host, d.dbCfg.Port, d.dbCfg.User, d.dbCfg.Password, d.dbCfg.Name, d.dbCfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Database) AutoMigrate(db *sql.DB, migrationsPath string) error {
	goose.SetDialect("postgres")
	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
