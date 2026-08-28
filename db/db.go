package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

func Init() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		os.Getenv("DB_PASSWORD"),
		host,
		port,
		os.Getenv("DB_NAME"),
	)

	log.Printf("Connecting to %v on port %v with username %v", host, port, user)
	db, err := gorm.Open(os.Getenv("DB_DRIVER"), dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused. AWS mariadb default is 8 hours.
	db.DB().SetConnMaxLifetime(time.Hour)

	return db
}

func Paginate(page, pageSize int) (int, int, func(db *gorm.DB) *gorm.DB) {
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 25
	}

	return page, pageSize, func(db *gorm.DB) *gorm.DB {
		offset := page * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
