package models

import "time"

type TyreModelFeature struct {
	ID          uint   `json:"id" gorm:"primary_key;auto_increment"`
	TyreModelID uint   `json:"tyre_model_id" gorm:"not null;unique_index"`
	FeatureJSON string `json:"-" gorm:"type:longtext;not null"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
