package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services/utils"
)

type Service struct {
	bannerStorage   Storage
	userBannerCache CacheStorage
}

type Storage interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (models.BannerContent, bool, error)
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

type CacheStorage interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (models.BannerContent, error)
	SaveUserBanner(ctx context.Context, tagID int, featureID int, value models.BannerContent)
}

func NewService(storage Storage, cacheStorage CacheStorage) *Service {
	return &Service{
		bannerStorage:   storage,
		userBannerCache: cacheStorage,
	}
}

var (
	ErrNotFound   = errors.New("banner not found")
	ErrNotAllowed = errors.New("not allowed")
)

func (s *Service) DeleteBanner(ctx context.Context, bannerID int) error {
	if !utils.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	affected, err := s.bannerStorage.DeleteBanner(ctx, bannerID)
	if err != nil {
		return err
	}

	if !affected {
		return ErrNotFound
	}

	return nil
}

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
