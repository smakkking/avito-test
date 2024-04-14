package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type BannerContent map[string]interface{}

func (a BannerContent) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *BannerContent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
