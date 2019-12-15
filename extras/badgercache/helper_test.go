package badgercache

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/rohanthewiz/logger"
	"github.com/rohanthewiz/serr"
)

func TestBadger(t *testing.T) {
	// create a new store with data under /tmp
	store, err := NewStore()
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error creating new key-value store"), true)
	}

	defer func() {
		err := store.Close()
		if err != nil {
			LogErr(err, "Error closing store - some values may not have been saved")
		}
	}()

	const hello_str = "hello"
	const world_str = "world"

	// Test Raw Bytes Set and Get
	key_byts := []byte(hello_str)
	test_val := []byte(world_str)
	err = store.SetBytes(key_byts, test_val) // save
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error setting key: "+hello_str), false)
	}

	byts, err := store.GetBytes(key_byts)
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error retrieving value for key: "+hello_str), false)
	}
	if string(byts) != string(test_val) {
		err = serr.Wrap(errors.New("Incorrect value retrieved from badger"), "retrieved", string(byts), "expected", world_str)
		handleTestErr(t, err, false)
	} else {
		fmt.Println("Yay! Retrieved the stored byte value:", string(byts), "for key:", hello_str)
	}

	// Test Strings Set and Get
	err = store.SetString(hello_str, world_str) // save
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error setting key: "+hello_str), false)
	}

	str, err := store.GetString(hello_str)
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error retrieving value for key: "+hello_str), false)
	}
	if str != world_str {
		err = serr.Wrap(errors.New("Incorrect value retrieved from badger"), "retrieved", str, "expected", world_str)
		handleTestErr(t, err, false)
	} else {
		fmt.Println("Yay! Retrieved the value:", str, "for key:", hello_str)
	}

	// Test Hashed Key with Touch
	const path_str = "/path/to/a/file.jpg"
	exists, err := store.ExistsHashed(path_str)
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error checking if key exists", "key", path_str), true)
	}
	if exists {
		handleTestErr(t, serr.Wrap(errors.New("Key should not already exist"), "key", path_str), true)
	}
	err = store.TouchHashed(path_str)
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error storing hash"), false)
	}
	exists_after_touch, err := store.ExistsHashed(path_str)
	if err != nil {
		handleTestErr(t, serr.Wrap(err, "Error on Exists()"), true)
	}
	if !exists_after_touch {
		handleTestErr(t, serr.Wrap(errors.New("key should exist after Touch()"), "key", path_str), false)
	}
}

func handleTestErr(t *testing.T, err error, fatal bool) {
	if err != nil {
		LogErr(err)
		if fatal {
			t.FailNow()
		}
	}
}
