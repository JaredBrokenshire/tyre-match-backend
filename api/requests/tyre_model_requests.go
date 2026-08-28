package requests

type CreateTyreModelRequest struct {
	Manufacturer      string `json:"manufacturer" example:"Michelin" validate:"required,max=200"`
	ModelName         string `json:"model_name" example:"Pilot Sport 5" validate:"required,max=200"`
	Category          string `json:"category" example:"Winter" validate:"max=200"`
	VehicleType       string `json:"vehicle_type" example:"Passenger Car" validate:"max=200"`
	WidthMm           int    `json:"width_mm" example:"205" validate:"min=0"`
	AspectRatio       int    `json:"aspect_ratio" example:"45" validate:"min=0"`
	RimDiameterInches int    `json:"rim_diameter_inches" example:"17" validate:"min=0"`
	GrooveCount       int    `json:"groove_count" example:"4" validate:"min=0"`
	PatternType       string `json:"pattern_type" example:"Symmetric" validate:"max=200"`
}

type UpdateTyreModelRequest struct {
	Manufacturer      string `json:"manufacturer" example:"Michelin" validate:"required,max=200"`
	ModelName         string `json:"model_name" example:"Pilot Sport 5" validate:"required,max=200"`
	Category          string `json:"category" example:"Winter" validate:"max=200"`
	VehicleType       string `json:"vehicle_type" example:"Passenger Car" validate:"max=200"`
	WidthMm           int    `json:"width_mm" example:"205" validate:"min=0"`
	AspectRatio       int    `json:"aspect_ratio" example:"45" validate:"min=0"`
	RimDiameterInches int    `json:"rim_diameter_inches" example:"17" validate:"min=0"`
	GrooveCount       int    `json:"groove_count" example:"4" validate:"min=0"`
	PatternType       string `json:"pattern_type" example:"Symmetric" validate:"max=200"`
}
