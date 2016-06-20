package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var dbPath = "./model-test"
var db *gorm.DB

func mkDB() {
	err := os.MkdirAll(dbPath, 0700)
	if err != nil {
		panic(err)
	}

	db, err = InitDB(fmt.Sprintf("%s/test.db", dbPath), false)
	if err != nil {
		panic(err)
	}
}

func rmDB() {
	if err := os.RemoveAll(dbPath); err != nil {
		panic(err)
	}
}
