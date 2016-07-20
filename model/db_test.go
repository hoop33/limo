package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

var dbPath = "./test.db"
var db *gorm.DB

func TestMain(m *testing.M) {
	rmDB()
	mkDB()
	run := m.Run()
	rmDB()
	os.Exit(run)
}

func mkDB() {
	var err error
	db, err = InitDB(dbPath, false)
	if err != nil {
		panic(err)
	}
}

func clearDB() {
	for _, table := range []string{
		"services",
		"stars",
		"tags",
		"star_tags",
	} {
		db.Exec(fmt.Sprintf("delete from %s", table))
	}
}

func rmDB() {
	if err := os.RemoveAll(dbPath); err != nil {
		panic(err)
	}
}
