package requests

type CreateTyreModelRequest struct {
	Manufacturer      string  `json:"manufacturer" example:"Michelin" validate:"required,max=200"`
	ModelName         string  `json:"model_name" example:"Pilot Sport 5" validate:"required,max=200"`
	WidthMm           int     `json:"width_mm" example:"205" validate:"min=0"`
	AspectRatio       int     `json:"aspect_ratio" example:"45" validate:"min=0"`
	RimDiameterInches int     `json:"rim_diameter_inches" example:"17" validate:"min=0"`
	GrooveCount       int     `json:"groove_count" example:"4" validate:"min=0"`
	PixelsPerInch     float32 `json:"pixels_per_inch" example:"200.0" validate:"min=0"`
	ROITop            int     `json:"roi_top" validate:"required,min=0"`
	ROILeft           int     `json:"roi_left" validate:"required,min=0"`
	ROIRight          int     `json:"roi_right" validate:"required,min=0"`
	ROIBottom         int     `json:"roi_bottom" validate:"required,min=0"`
}

type UpdateTyreModelRequest struct {
	Manufacturer      string `json:"manufacturer" example:"Michelin" validate:"required,max=200"`
	ModelName         string `json:"model_name" example:"Pilot Sport 5" validate:"required,max=200"`
	WidthMm           int    `json:"width_mm" example:"205" validate:"min=0"`
	AspectRatio       int    `json:"aspect_ratio" example:"45" validate:"min=0"`
	RimDiameterInches int    `json:"rim_diameter_inches" example:"17" validate:"min=0"`
	GrooveCount       int    `json:"groove_count" example:"4" validate:"min=0"`
}
