package factories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
	"log"
	"tyre-match-backend/db/models"
)

func NewFile(db *gorm.DB, file *models.File) {
	fillFileDefaults(file)
	if err := db.Create(&file).Error; err != nil {
		log.Println("Error creating file in factory: ", err.Error())
	}
}

func fillFileDefaults(file *models.File) {
	if file.Model == "" {
		file.Model = models.FileModelTyreModel
	}
	if file.ModelId == 0 {
		file.ModelId = 1
	}
	if file.FileType == "" {
		file.FileType = models.FileTypeOriginal
	}
	if file.Name == "" {
		file.Name = random.String(16) + ".jpg"
	}
	if file.Location == "" {
		file.Location = fmt.Sprintf("/%v", random.String(16))
	}
}
