package responses

import (
	"time"
	m "tyre-match-backend/db/models"
)

type TyreModelResponse struct {
	ID                uint   `json:"id"`
	Manufacturer      string `json:"manufacturer"`
	ModelName         string `json:"model_name"`
	WidthMm           int    `json:"width_mm"`
	AspectRatio       int    `json:"aspect_ratio"`
	RimDiameterInches int    `json:"rim_diameter_inches"`
	GrooveCount       int    `json:"groove_count"`

	Status        string  `json:"status"`
	PixelsPerInch float32 `json:"pixels_per_inch"`

	ROITop    int `json:"roi_top"`
	ROILeft   int `json:"roi_left"`
	ROIRight  int `json:"roi_right"`
	ROIBottom int `json:"roi_bottom"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Images map[string]FileResponse `json:"images"`
}

type SlimTyreModelResponse struct {
	ID           uint   `json:"id"`
	Manufacturer string `json:"manufacturer"`
	ModelName    string `json:"model_name"`
	Status       string `json:"status"`

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
		WidthMm:           tyreModel.WidthMm,
		AspectRatio:       tyreModel.AspectRatio,
		RimDiameterInches: tyreModel.RimDiameterInches,
		GrooveCount:       tyreModel.GrooveCount,
		Status:            tyreModel.Status,
		PixelsPerInch:     tyreModel.PixelsPerInch,
		ROITop:            tyreModel.ROITop,
		ROILeft:           tyreModel.ROILeft,
		ROIRight:          tyreModel.ROIRight,
		ROIBottom:         tyreModel.ROIBottom,
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
		Status:       tyreModel.Status,
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
