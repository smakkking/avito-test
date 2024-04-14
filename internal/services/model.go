package services

import (
	"context"
	"errors"

	"github.com/smakkking/avito_test/internal/models"
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
