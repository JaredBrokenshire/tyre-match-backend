package models

import "time"

type TyreModel struct {
	ID uint `gorm:"primary_key;auto_increment" json:"id"`

	Manufacturer string `json:"manufacturer"`
	ModelName    string `json:"model_name"`

	WidthMm           int `json:"width_mm"`
	AspectRatio       int `json:"aspect_ratio"`
	RimDiameterInches int `json:"rim_diameter_inches"`
	GrooveCount       int `json:"groove_count"`

	Status        string  `json:"status"`
	PixelsPerInch float32 `json:"pixels_per_inch"`

	ROITop    int `json:"roi_top"`
	ROILeft   int `json:"roi_left"`
	ROIRight  int `json:"roi_right"`
	ROIBottom int `json:"roi_bottom"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Images map[string]*File `json:"images" gorm:"-"`
}
