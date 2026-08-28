package models

import "time"

type TyreModel struct {
	ID uint `gorm:"primary_key;auto_increment" json:"id"`

	Manufacturer string `json:"manufacturer"`
	ModelName    string `json:"model_name"`

	Category    string `json:"category"`
	VehicleType string `json:"vehicle_type"`

	WidthMm           int    `json:"width_mm"`
	AspectRatio       int    `json:"aspect_ratio"`
	RimDiameterInches int    `json:"rim_diameter_inches"`
	GrooveCount       int    `json:"groove_count"`
	PatternType       string `json:"pattern_type"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Images map[string]*File `json:"images" gorm:"-"`
}
