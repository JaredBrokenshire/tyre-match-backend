package repositories

import (
	"github.com/jinzhu/gorm"
	"tyre-match-backend/db"
	m "tyre-match-backend/db/models"
)

type TyreModelRepository struct {
	*Repository
}

func NewTyreModelRepository(db *gorm.DB) *TyreModelRepository {
	return &TyreModelRepository{
		Repository: &Repository{
			Db: db,
		},
	}
}

func (r *TyreModelRepository) List(page, pageSize int, scopes []func(db *gorm.DB) *gorm.DB) ([]*m.TyreModel, int, int, int) {
	var tyreModels []*m.TyreModel
	var count int64

	page, pageSize, paginateFunc := db.Paginate(page, pageSize)

	r.Db.
		Scopes(paginateFunc).
		Scopes(scopes...).
		Order("model_name ASC").
		Find(&tyreModels)

	// Separate count to include records filtered out by pagination
	r.Db.
		Model(&m.TyreModel{}).
		Scopes(scopes...).
		Count(&count)

	return tyreModels, int(count), page, pageSize
}

func (r *TyreModelRepository) GetByID(id uint) *m.TyreModel {
	var tyreModel m.TyreModel
	r.Db.Where("id = ?", id).First(&tyreModel)
	if tyreModel.ID == 0 {
		return nil
	}

	tyreModel.Images = make(map[string]*m.File)

	// Load files
	var images []m.File
	r.Db.Where("model = ?", m.FileModelTyreModel).Where("model_id = ?", id).Find(&images)
	for i := range images {
		file := &images[i]
		tyreModel.Images[file.FileType] = file
	}

	return &tyreModel
}
