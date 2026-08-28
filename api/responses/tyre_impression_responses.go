package responses

import (
	"time"
	m "tyre-match-backend/db/models"
)

type TyreImpressionResponse struct {
	ID            uint    `json:"id"`
	Status        string  `json:"status"`
	PixelsPerInch float32 `json:"pixels_per_inch"`
	EdgeDensity   float32 `json:"edge_density"`
	VoidRatio     float32 `json:"void_ratio"`
	GrooveCount   int     `json:"groove_count"`

	Images map[string]FileResponse `json:"images"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SlimTyreImpressionResponse struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TyreImpressionPaginatedResponse struct {
	Data []SlimTyreImpressionResponse `json:"data"`
	Meta ResponseMeta                 `json:"meta"`
}

func NewTyreImpressionResponse(tyreImpression *m.TyreImpression) *TyreImpressionResponse {
	res := &TyreImpressionResponse{
		ID:     tyreImpression.ID,
		Status: tyreImpression.Status,

		PixelsPerInch: tyreImpression.PixelsPerInch,
		EdgeDensity:   tyreImpression.EdgeDensity,
		VoidRatio:     tyreImpression.VoidRatio,
		GrooveCount:   tyreImpression.GrooveCount,

		CreatedAt: tyreImpression.CreatedAt,
		UpdatedAt: tyreImpression.UpdatedAt,
	}

	res.Images = make(map[string]FileResponse)
	for key, value := range tyreImpression.Images {
		res.Images[key] = *NewFileResponse(value)
	}

	return res
}

func NewSlimTyreImpressionResponse(tyreImpression *m.TyreImpression) *SlimTyreImpressionResponse {
	return &SlimTyreImpressionResponse{
		ID:        tyreImpression.ID,
		Status:    tyreImpression.Status,
		CreatedAt: tyreImpression.CreatedAt,
		UpdatedAt: tyreImpression.UpdatedAt,
	}
}

func NewTyreImpressionResponses(tyreImpressions []*m.TyreImpression) []SlimTyreImpressionResponse {
	var res []SlimTyreImpressionResponse
	for _, tyreImpression := range tyreImpressions {
		res = append(res, *NewSlimTyreImpressionResponse(tyreImpression))
	}
	return res
}

func NewTyreImpressionPaginatedResponse(tyreImpressions []*m.TyreImpression, count, page, pageSize int) *TyreImpressionPaginatedResponse {
	return &TyreImpressionPaginatedResponse{
		Data: NewTyreImpressionResponses(tyreImpressions),
		Meta: ResponseMeta{
			TotalCount: count,
			Page:       page,
			PageSize:   pageSize,
		},
	}
}
