package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services/utils"
)

func (s *Service) GetUserBanner(ctx context.Context, tagID, featureID int, useLastRevision bool) (models.BannerContent, error) {
	if !useLastRevision {
		banner, err := s.userBannerCache.GetUserBanner(ctx, tagID, featureID)
		if err == nil {
			return banner, nil
		}
	}

	banner, enabled, err := s.bannerStorage.GetUserBanner(ctx, tagID, featureID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if !utils.IsAdmin(ctx) && !enabled {
		return nil, ErrNotAllowed
	}

	s.userBannerCache.SaveUserBanner(ctx, tagID, featureID, banner)
	return banner, nil
}
