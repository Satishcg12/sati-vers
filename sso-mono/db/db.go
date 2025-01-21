package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/Satishcg12/sati-vers/sso-mono/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Database struct {
	dbCfg config.Database
}

type IDatabase interface {
	Connect() (*sql.DB, error)
	AutoMigrate(*sql.DB, string) error
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

func (d *Database) AutoMigrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")
	err := goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
