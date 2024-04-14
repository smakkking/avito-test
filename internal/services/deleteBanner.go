package services

import (
	"context"

	"github.com/smakkking/avito_test/internal/services/utils"
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
