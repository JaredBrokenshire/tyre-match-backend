package responses

import m "tyre-match-backend/db/models"

type TyreModelFeatureResponse struct {
	FeatureJSON string `json:"feature_json"`
}

func NewTyreModelFeatureResponse(tyreModelFeature *m.TyreModelFeature) *TyreModelFeatureResponse {
	return &TyreModelFeatureResponse{
		FeatureJSON: tyreModelFeature.FeatureJSON,
	}
}
