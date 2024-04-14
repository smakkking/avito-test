package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
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

func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) (bool, error) {
	err := s.db.QueryRowContext(
		ctx,
		`
		DELETE FROM BannersInfo
		WHERE "id" = $1
		RETURNING "id"
		`, bannerID,
	).Scan(&bannerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Storage) UpdateBanner(ctx context.Context, bannerID int, banner *models.BasicBannnerInfo) (bool, error) {
	err := s.db.QueryRowContext(
		ctx,
		`
		UPDATE BannersInfo
		SET "value" = $1, tag_array = $2, feature = $3, is_enabled = $4, updated_at = now()
		WHERE "id" = $5
		RETURNING "id"
		`,
		banner.Content, pq.Array(banner.TagIDs), banner.FeatureID, banner.IsActive, bannerID,
	).Scan(&bannerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Storage) CreateBanner(ctx context.Context, banner *models.BasicBannnerInfo) (int, error) {
	var bannerID int

	err := s.db.QueryRowContext(
		ctx,
		`
		INSERT INTO BannersInfo("value", tag_array, feature, is_enabled)
		VALUES ($1::jsonb, $2, $3, $4) RETURNING "id";
		`,
		banner.Content, pq.Array(banner.TagIDs), banner.FeatureID, banner.IsActive,
	).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (s *Storage) GetUserBanner(ctx context.Context, tagID int, featureID int) (models.BannerContent, bool, error) {
	var data models.BannerContent
	var isEnabled bool

	err := s.db.QueryRowContext(
		ctx,
		`SELECT "value", is_enabled 
		   FROM BannersInfo 
		  WHERE feature = $1 AND array_position(tag_array, $2) IS NOT NULL;`,
		featureID, tagID,
	).Scan(&data, &isEnabled)
	if err != nil {
		return nil, false, err
	}

	return data, isEnabled, nil
}

func (s *Storage) GetAllBannersFiltered(
	ctx context.Context,
	tagID int, tagSearch bool,
	featureID int, featureSearch bool,
	limit int, offset int,
) ([]*models.BannerInfo, error) {
	query := strings.Builder{}
	queryArgs := make([]interface{}, 0)

	query.WriteString(`SELECT "id", tag_array, "value", feature, is_enabled, created_at, updated_at FROM BannersInfo `)

	// фильтрация
	query.WriteString(`WHERE `)

	if tagSearch {
		query.WriteString(`array_position(tag_array, $` + fmt.Sprint(len(queryArgs)+1) + `) IS NOT NULL AND `)
		queryArgs = append(queryArgs, tagID)
	}

	if featureSearch {
		query.WriteString(`feature = $` + fmt.Sprint(len(queryArgs)+1) + ` AND`)
		queryArgs = append(queryArgs, featureID)
	}

	query.WriteString(`true `)

	// смещение
	if offset != -1 {
		query.WriteString(`OFFSET $` + fmt.Sprint(len(queryArgs)+1) + ` `)
		queryArgs = append(queryArgs, offset)
	}

	if limit != -1 {
		query.WriteString(`LIMIT $` + fmt.Sprint(len(queryArgs)+1) + ` `)
		queryArgs = append(queryArgs, limit)
	}

	rows, err := s.db.QueryContext(
		ctx,
		query.String(),
		queryArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.BannerInfo, 0)

	for rows.Next() {
		banner := new(models.BannerInfo)

		err := rows.Scan(
			&banner.BannerID,
			pq.Array(&banner.TagIDs), // почитать в статье
			&banner.Content,
			&banner.FeatureID,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, banner)
	}

	return result, nil
}
