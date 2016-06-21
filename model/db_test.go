package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

var dbPath = "./model-test"
var db *gorm.DB

func TestMain(m *testing.M) {
	rmDB()
	mkDB()
	run := m.Run()
	rmDB()
	os.Exit(run)
}

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
