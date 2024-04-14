package postgres

import (
	"context"

	"github.com/smakkking/avito_test/internal/models"
)

func (s *Storage) GetUserBanner(ctx context.Context, tagID, featureID int) (models.BannerContent, bool, error) {
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
