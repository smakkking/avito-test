package models

import "time"

type BannerInfo struct {
	BannerID  int         `json:"banner_id"`
	TagIDs    []int       `json:"tag_ids"`
	FeatureID int         `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
