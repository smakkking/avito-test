package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

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
