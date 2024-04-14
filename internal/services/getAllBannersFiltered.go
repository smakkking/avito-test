package services

import (
	"context"

	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services/utils"
)

func (s *Service) GetAllBannersFiltered(
	ctx context.Context,
	tagID int, tagSearch bool,
	featureID int, featureSearch bool,
	limit int, offset int,
) ([]*models.BannerInfo, error) {
	if !utils.IsAdmin(ctx) {
		return nil, ErrNotAllowed
	}

	banners, err := s.bannerStorage.GetAllBannersFiltered(ctx, tagID, tagSearch, featureID, featureSearch, limit, offset)
	if err != nil {
		return nil, err
	}

	return banners, err
}
