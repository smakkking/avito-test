package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/smakkking/avito_test/internal/app"
	"github.com/smakkking/avito_test/internal/models"
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

func (s *Storage) GetUserBanner(ctx context.Context, tagID int, featureID int) (interface{}, bool, error) {
	var data interface{}
	var isEnabled bool

	err := s.db.QueryRowContext(
		ctx,
		`SELECT "value", is_enabled 
		   FROM BannersInfo 
		  WHERE tag = $1 AND feature = $2`,
		tagID, featureID,
	).Scan(&data, &isEnabled)
	if err != nil {
		return nil, false, err
	}

	return data, isEnabled, nil
}

func (s *Storage) GetAllBannersFiltered(ctx context.Context, tagID int, featureID int, limit int, offset int) ([]models.BannerInfo, error) {
	// сортировать по возрастанию id

	rows, err := s.db.QueryContext(
		ctx, 
		`SELECT DISTINCT b_id FROM BannersInfo WHERE `
	)
}
