package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/smakkking/avito_test/internal/models"
)

func (s *Storage) GetAllBannersFiltered(
	ctx context.Context,
	tagID int, tagSearch bool,
	featureID int, featureSearch bool,
	limit int, offset int,
) ([]*models.BannerInfo, error) {
	query := strings.Builder{}
	queryArgs := make([]interface{}, 0)

	query.WriteString(`SELECT "id", tag_array, "value", feature, is_enabled, created_at, updated_at FROM BannersInfo `)

	// фильтрация
	query.WriteString(`WHERE `)

	if tagSearch {
		query.WriteString(`array_position(tag_array, $` + fmt.Sprint(len(queryArgs)+1) + `) IS NOT NULL AND `)
		queryArgs = append(queryArgs, tagID)
	}

	if featureSearch {
		query.WriteString(`feature = $` + fmt.Sprint(len(queryArgs)+1) + ` AND`)
		queryArgs = append(queryArgs, featureID)
	}

	query.WriteString(`true `)

	// смещение
	if offset != -1 {
		query.WriteString(`OFFSET $` + fmt.Sprint(len(queryArgs)+1) + ` `)
		queryArgs = append(queryArgs, offset)
	}

	if limit != -1 {
		query.WriteString(`LIMIT $` + fmt.Sprint(len(queryArgs)+1) + ` `)
		queryArgs = append(queryArgs, limit)
	}

	rows, err := s.db.QueryContext(
		ctx,
		query.String(),
		queryArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.BannerInfo, 0)

	for rows.Next() {
		banner := new(models.BannerInfo)

		x := make([]sql.NullInt32, 0)

		err := rows.Scan(
			&banner.BannerID,
			pq.Array(&x),
			&banner.Content,
			&banner.FeatureID,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		for _, val := range x {
			banner.TagIDs = append(banner.TagIDs, int(val.Int32))
		}

		result = append(result, banner)
	}

	return result, nil
}
