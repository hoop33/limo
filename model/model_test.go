package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// Use the sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbPath = "./model-test"
var db *gorm.DB

func mkDB() {
	err := os.MkdirAll(dbPath, 0700)
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open("sqlite3", fmt.Sprintf("%s/test.db", dbPath))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Service{}, &Star{}, &Tag{})
}

func rmDB() {
	if err := os.RemoveAll(dbPath); err != nil {
		panic(err)
	}
}
