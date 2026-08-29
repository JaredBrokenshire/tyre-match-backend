package models

import "time"

const (
	FileModelTyreImpression = "TyreImpression"
	FileModelTyreModel      = "TyreModel"
)

const (
	FileTypeOriginal   = "original"
	FileTypeNormalised = "normalised"
	FileTypeEnhanced   = "enhanced"
	FileTypeBinary     = "binary"
)

type File struct {
	ID       uint   `json:"-" gorm:"primary_key"`
	Model    string `json:"model"`
	ModelId  uint   `json:"model_id"`
	FileType string `json:"file_type"`

	Name     string `json:"name"`
	Location string `json:"location"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
