package services

import (
	"context"
	"errors"

	"github.com/smakkking/avito_test/internal/models"
)

type Service struct {
	bannerStorage Storage
}

type Storage interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (interface{}, bool, error)
	GetAllBannersFiltered(
		ctx context.Context,
		tagID int, tagSearch bool,
		featureID int, featureSearch bool,
		limit int, offset int,
	) ([]*models.BannerInfo, error)
	CreateBanner(ctx context.Context, banner *models.BasicBannnerInfo) (int, error)
	UpdateBanner(ctx context.Context, bannerID int, banner *models.BasicBannnerInfo) (bool, error)
	DeleteBanner(ctx context.Context, bannerID int) (bool, error)
}

func NewService(storage Storage) *Service {
	return &Service{
		bannerStorage: storage,
	}
}

var (
	ErrNotFound = errors.New("banner not found")
)

func (s *Service) DeleteBanner(ctx context.Context, bannerID int) (bool, error) {
	return s.bannerStorage.DeleteBanner(ctx, bannerID)
}

func (s *Service) UpdateBanner(ctx context.Context, bannerID int, banner *models.BasicBannnerInfo) (bool, error) {
	return s.bannerStorage.UpdateBanner(ctx, bannerID, banner)
}

func (s *Service) GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool) (interface{}, error) {
	banner, _, err := s.bannerStorage.GetUserBanner(ctx, tagID, featureID)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (s *Service) GetAllBannersFiltered(
	ctx context.Context,
	tagID int, tagSearch bool,
	featureID int, featureSearch bool,
	limit int, offset int,
) ([]*models.BannerInfo, error) {
	banners, err := s.bannerStorage.GetAllBannersFiltered(ctx, tagID, tagSearch, featureID, featureSearch, limit, offset)
	if err != nil {
		return nil, err
	}

	return banners, err
}

func (s *Service) CreateBanner(ctx context.Context, banner *models.BasicBannnerInfo) (int, error) {
	bannerID, err := s.bannerStorage.CreateBanner(ctx, banner)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}
