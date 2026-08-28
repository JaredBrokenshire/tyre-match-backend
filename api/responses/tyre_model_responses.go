package responses

import (
	"time"
	m "tyre-match-backend/db/models"
)

type TyreModelResponse struct {
	ID                uint   `json:"id"`
	Manufacturer      string `json:"manufacturer"`
	ModelName         string `json:"model_name"`
	Category          string `json:"category"`
	VehicleType       string `json:"vehicle_type"`
	WidthMm           int    `json:"width_mm"`
	AspectRatio       int    `json:"aspect_ratio"`
	RimDiameterInches int    `json:"rim_diameter_inches"`
	GrooveCount       int    `json:"groove_count"`
	PatternType       string `json:"pattern_type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Images map[string]FileResponse `json:"images"`
}

type SlimTyreModelResponse struct {
	ID           uint   `json:"id"`
	Manufacturer string `json:"manufacturer"`
	ModelName    string `json:"model_name"`
	Category     string `json:"category"`
	VehicleType  string `json:"vehicle_type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TyreModelPaginatedResponse struct {
	Data []SlimTyreModelResponse `json:"data"`
	Meta ResponseMeta            `json:"meta"`
}

func NewTyreModelResponse(tyreModel *m.TyreModel) *TyreModelResponse {
	res := &TyreModelResponse{
		ID:                tyreModel.ID,
		Manufacturer:      tyreModel.Manufacturer,
		ModelName:         tyreModel.ModelName,
		Category:          tyreModel.Category,
		VehicleType:       tyreModel.VehicleType,
		WidthMm:           tyreModel.WidthMm,
		AspectRatio:       tyreModel.AspectRatio,
		RimDiameterInches: tyreModel.RimDiameterInches,
		GrooveCount:       tyreModel.GrooveCount,
		PatternType:       tyreModel.PatternType,
		CreatedAt:         tyreModel.CreatedAt,
		UpdatedAt:         tyreModel.UpdatedAt,
	}

	res.Images = make(map[string]FileResponse)
	for key, value := range tyreModel.Images {
		res.Images[key] = *NewFileResponse(value)
	}

	return res
}

func NewSlimTyreModelResponse(tyreModel *m.TyreModel) *SlimTyreModelResponse {
	return &SlimTyreModelResponse{
		ID:           tyreModel.ID,
		Manufacturer: tyreModel.Manufacturer,
		ModelName:    tyreModel.ModelName,
		Category:     tyreModel.Category,
		VehicleType:  tyreModel.VehicleType,
		CreatedAt:    tyreModel.CreatedAt,
		UpdatedAt:    tyreModel.UpdatedAt,
	}
}

func NewTyreModelResponses(tyreModels []*m.TyreModel) []SlimTyreModelResponse {
	var res []SlimTyreModelResponse
	for _, tyreModel := range tyreModels {
		res = append(res, *NewSlimTyreModelResponse(tyreModel))
	}
	return res
}

func NewTyreModelPaginatedResponse(tyreModels []*m.TyreModel, count, page, pageSize int) *TyreModelPaginatedResponse {
	return &TyreModelPaginatedResponse{
		Data: NewTyreModelResponses(tyreModels),
		Meta: ResponseMeta{
			TotalCount: count,
			Page:       page,
			PageSize:   pageSize,
		},
	}
}
