package repositories

import (
	"github.com/jinzhu/gorm"
	m "tyre-match-backend/db/models"
)

type TyreModelFeatureRepository struct {
	*Repository
}

func NewTyreModelFeatureRepository(db *gorm.DB) *TyreModelFeatureRepository {
	return &TyreModelFeatureRepository{Repository: &Repository{Db: db}}
}

func (r *TyreModelFeatureRepository) GetByTyreModelID(id uint) *m.TyreModelFeature {
	var feature m.TyreModelFeature
	r.Db.Where("tyre_model_id = ?", id).First(&feature)
	if feature.ID == 0 {
		return nil
	}
	return &feature
}

func (r *TyreModelFeatureRepository) List() []*m.TyreModelFeature {
	var features []*m.TyreModelFeature
	r.Db.Order("tyre_model_id ASC").Find(&features)
	return features
}

func (r *TyreModelFeatureRepository) Upsert(feature *m.TyreModelFeature) error {
	if feature == nil {
		return gorm.ErrInvalidSQL
	}
	existing := r.GetByTyreModelID(feature.TyreModelID)
	if existing == nil {
		return r.Create(feature)
	}
	existing.FeatureJSON = feature.FeatureJSON
	return r.Update(existing)
}
