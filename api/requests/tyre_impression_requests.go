package requests

type CreateTyreImpressionRequest struct {
	PixelsPerInch float32 `json:"pixels_per_inch" validate:"required,min=0"`
	ROITop        int     `json:"roi_top" validate:"required,min=0"`
	ROILeft       int     `json:"roi_left" validate:"required,min=0"`
	ROIRight      int     `json:"roi_right" validate:"required,gt=0"`
	ROIBottom     int     `json:"roi_bottom" validate:"required,gt=0"`
}
