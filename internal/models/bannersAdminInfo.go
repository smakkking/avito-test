package models

import (
	"time"
)

type BasicBannnerInfo struct {
	TagIDs    []int         `json:"tag_ids"`
	FeatureID int           `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  bool          `json:"is_active"`
}

type BannerInfo struct {
	BasicBannnerInfo
	BannerID  int       `json:"banner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
