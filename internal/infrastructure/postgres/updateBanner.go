package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/smakkking/avito_test/internal/models"
)

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
