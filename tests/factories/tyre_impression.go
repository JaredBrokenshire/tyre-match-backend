package factories

import (
	"github.com/jinzhu/gorm"
	"log"
	m "tyre-match-backend/db/models"
)

func NewTyreImpression(db *gorm.DB, tyreImpression *m.TyreImpression) {
	fillTyreImpressionDefaults(tyreImpression)
	err := db.Create(tyreImpression).Error
	if err != nil {
		log.Println("Error creating tyre impression in factory: ", err.Error())
	}
}

func fillTyreImpressionDefaults(tyreImpression *m.TyreImpression) {
	if tyreImpression.PixelsPerInch == 0 {
		tyreImpression.PixelsPerInch = 1
	}
	if tyreImpression.Status == "" {
		tyreImpression.Status = m.ProcessingStatusUploaded
	}
}
