package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/smakkking/avito_test/internal/app"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfg app.Config) (*Storage, error) {
	time.Sleep(5 * time.Second) // для корректного подключения в докере

	database_url := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDBName,
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgSSLMode,
	)
	db, err := sql.Open("postgres", database_url)
	if err != nil {
		return nil, err
	}

	// проверка, что подключились
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
