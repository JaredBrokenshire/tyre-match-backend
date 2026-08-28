package repositories

import "github.com/jinzhu/gorm"

type Repos struct {
	DB             *gorm.DB
	TyreModel      *TyreModelRepository
	TyreImpression *TyreImpressionRepository
	File           *FileRepository
}

func NewRepos(db *gorm.DB) *Repos {
	return &Repos{
		DB:             db,
		TyreModel:      NewTyreModelRepository(db),
		TyreImpression: NewTyreImpressionRepository(db),
		File:           NewFileRepository(db),
	}
}
