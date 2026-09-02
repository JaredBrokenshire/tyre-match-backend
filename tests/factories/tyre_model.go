package factories

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
	"log"
	m "tyre-match-backend/db/models"
)

func NewTyreModel(db *gorm.DB, tyreModel *m.TyreModel) {
	fillTyreModelDefaults(tyreModel)
	err := db.Create(tyreModel).Error
	if err != nil {
		log.Println("Error creating tyre model in factory: ", err.Error())
	}
}

func fillTyreModelDefaults(tyreModel *m.TyreModel) {
	if tyreModel.Manufacturer == "" {
		tyreModel.Manufacturer = random.String(16)
	}
	if tyreModel.ModelName == "" {
		tyreModel.ModelName = random.String(16)
	}
	if tyreModel.Status == "" {
		tyreModel.Status = m.ProcessingStatusProcessed
	}
}
