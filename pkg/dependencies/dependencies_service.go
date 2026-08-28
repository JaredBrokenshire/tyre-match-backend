package dependencies

import (
	"github.com/jinzhu/gorm"
	"tyre-match-backend/pkg/file_storage"
)

type DependencyService struct {
	db *gorm.DB

	fileStore file_storage.Store
}

func NewDependencyService(db *gorm.DB) *DependencyService {
	ds := &DependencyService{db: db}
	return ds
}

func (s *DependencyService) GetDB() *gorm.DB {
	return s.db
}

func (s *DependencyService) GetFileStore() file_storage.Store {
	if s.fileStore == nil {
		s.fileStore = file_storage.NewLocalStorage("assets/files/")
	}

	return s.fileStore
}

func (s *DependencyService) SetFileStore(store file_storage.Store) {
	s.fileStore = store
}
