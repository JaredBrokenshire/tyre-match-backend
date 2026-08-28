package repositories

import "github.com/jinzhu/gorm"

type FileRepository struct {
	*Repository
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		Repository: &Repository{
			Db: db,
		},
	}
}
