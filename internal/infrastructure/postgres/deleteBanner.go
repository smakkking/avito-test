package postgres

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) (bool, error) {
	err := s.db.QueryRowContext(
		ctx,
		`
		DELETE FROM BannersInfo
		WHERE "id" = $1
		RETURNING "id"
		`, bannerID,
	).Scan(&bannerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
