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
	GetAllBannersFiltered(ctx context.Context, tagID int, featureID int, limit int, offset int) ([]models.BannerInfo, error)
}

func NewService(storage Storage) *Service {
	return &Service{
		bannerStorage: storage,
	}
}

var (
	ErrNotFound = errors.New("banner not found")
)

func (s *Service) GetUserBanner(ctx context.Context, tagID int, featureID int, useLastRevision bool) (interface{}, error) {
	banner, _, err := s.bannerStorage.GetUserBanner(ctx, tagID, featureID)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (s *Service) GetAllBannersFiltered(ctx context.Context, tagID int, featureID int, limit int, offset int) ([]models.BannerInfo, error) {
	banners, err := s.bannerStorage.GetAllBannersFiltered(ctx, tagID, featureID, limit, offset)
	if err != nil {
		return nil, err
	}

	return banners, err
}
