package postgres

import (
	"context"

	"github.com/lib/pq"
	"github.com/smakkking/avito_test/internal/models"
)

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
