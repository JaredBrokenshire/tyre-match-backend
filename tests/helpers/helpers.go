package helpers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

func ClearTable(db *gorm.DB, tableName string) {
	err := db.Exec(fmt.Sprintf("DELETE FROM %v", tableName)).Error
	if err != nil {
		log.Fatalf("You can't clear that table. Err: %v", err)
	}
	err = db.Exec(fmt.Sprintf("ALTER TABLE %v AUTO_INCREMENT = 1", tableName)).Error
	if err != nil {
		log.Fatalf("Error setting autoincrement. Err: %v", err)
	}
}
