package storage

import (
	"log"
	"testing"

	str "github.com/vestlog/nix/pkg/structs"
)

var (
	db  *GormDatabase
	dsn = "storage.db"
)

func create() {
	var err error
	db, err = CreateGormDatabase(dsn)
	if err != nil {
		log.Fatal("Could not create DB")
	}
}

func TestSavePost(t *testing.T) {
	if db == nil {
		create()
	}
	if err := db.db.AutoMigrate(&str.Post{}); err != nil {
		t.Error("TestSavePost failed:", err)
	}
	if err := db.SavePost(&str.Post{
		UserID: 0,
		ID:     0,
		Title:  "",
		Body:   "",
	}); err != nil {
		t.Error("TestSavePost failed:", err)
	}
}
