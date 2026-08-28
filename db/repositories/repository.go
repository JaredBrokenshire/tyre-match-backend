package repositories

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) Create(object interface{}) error {
	return r.Db.Create(object).Error
}

func (r *Repository) Update(object interface{}) error {
	return r.Db.Save(object).Error
}

func (r *Repository) Delete(object interface{}) error {
	return r.Db.Delete(object).Error
}

// Restore an object by setting deleted_at to nil
func (r *Repository) Restore(object interface{}) error {
	return r.Db.Unscoped().Model(object).Update("deleted_at", nil).Error
}

type Scopes []func(db *gorm.DB) *gorm.DB
