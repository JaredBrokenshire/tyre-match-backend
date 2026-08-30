package models

import "time"

const (
	TyreImpressionStatusUploaded   = "Uploaded"
	TyreImpressionStatusProcessing = "Processing"
	TyreImpressionStatusProcessed  = "Processed"
	TyreImpressionStatusMatched    = "Matched"
	TyreImpressionStatusFailed     = "Failed"
)

type TyreImpression struct {
	ID uint `json:"id" gorm:"primary_key;auto_increment"`

	Status string `json:"status"`

	PixelsPerInch float32 `json:"pixels_per_inch"`

	ROITop    int `json:"roi_top"`
	ROILeft   int `json:"roi_left"`
	ROIRight  int `json:"roi_right"`
	ROIBottom int `json:"roi_bottom"`

	EdgeDensity float32 `json:"edge_density"`
	VoidRatio     float32 `json:"void_ratio"`
	GrooveCount   int     `json:"groove_count"`

	Images map[string]*File `json:"images" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
