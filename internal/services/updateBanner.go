package services

import (
	"context"

	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services/utils"
)

func (s *Service) UpdateBanner(ctx context.Context, bannerID int, banner *models.BasicBannnerInfo) error {
	if !utils.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	affected, err := s.bannerStorage.UpdateBanner(ctx, bannerID, banner)
	if err != nil {
		return err
	}

	if !affected {
		return ErrNotFound
	}

	return nil
}
