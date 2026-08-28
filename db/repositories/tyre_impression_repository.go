package repositories

import (
	"github.com/jinzhu/gorm"
	"tyre-match-backend/db"
	m "tyre-match-backend/db/models"
)

type TyreImpressionRepository struct {
	*Repository
}

func NewTyreImpressionRepository(db *gorm.DB) *TyreImpressionRepository {
	return &TyreImpressionRepository{
		Repository: &Repository{
			Db: db,
		},
	}
}

func (r *TyreImpressionRepository) List(page, pageSize int, scopes []func(db *gorm.DB) *gorm.DB) ([]*m.TyreImpression, int, int, int) {
	var tyreImpressions []*m.TyreImpression
	var count int64

	page, pageSize, paginateFunc := db.Paginate(page, pageSize)

	r.Db.
		Scopes(paginateFunc).
		Scopes(scopes...).
		Order("id ASC").
		Find(&tyreImpressions)

	// Separate count to include records filtered out by pagination
	r.Db.
		Model(&m.TyreImpression{}).
		Scopes(scopes...).
		Count(&count)

	return tyreImpressions, int(count), page, pageSize
}

func (r *TyreImpressionRepository) GetByID(id uint) *m.TyreImpression {
	var tyreImpression m.TyreImpression
	r.Db.Where("id = ?", id).First(&tyreImpression)
	if tyreImpression.ID == 0 {
		return nil
	}

	tyreImpression.Images = make(map[string]*m.File)

	// Load files
	var images []m.File
	r.Db.Where("model = ?", m.FileModelTyreImpression).Where("model_id = ?", id).Find(&images)
	for i := range images {
		file := &images[i]
		tyreImpression.Images[file.FileType] = file
	}

	return &tyreImpression
}
