package services

import (
	"context"

	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services/utils"
)

func (s *Service) CreateBanner(ctx context.Context, banner *models.BasicBannnerInfo) (int, error) {
	if !utils.IsAdmin(ctx) {
		return 0, ErrNotAllowed
	}

	bannerID, err := s.bannerStorage.CreateBanner(ctx, banner)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}
