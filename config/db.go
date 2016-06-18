package config

import (
	"github.com/hoop33/limo/model"
	"github.com/jinzhu/gorm"
	// Use the sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// InitDB initializes the database at the specified path
func InitDB(filepath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Service{}, &model.Star{})

	return db, nil
}
