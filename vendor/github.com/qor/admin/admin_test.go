package admin

import (
	"testing"

	"github.com/qor/qor"
	"github.com/qor/qor/test/utils"
)

type User struct {
	Name string
	ID   uint64
}

var db = utils.TestDB()

func TestAddResource(t *testing.T) {
	admin := New(&qor.Config{DB: db})
	user := admin.AddResource(&User{})

	if user != admin.resources[0] {
		t.Error("resource not added")
	}

	if admin.GetMenus()[0].Name != "Users" {
		t.Error("resource not added to menu")
	}
}

func TestAddResourceWithInvisibleOption(t *testing.T) {
	admin := New(&qor.Config{DB: db})
	user := admin.AddResource(&User{}, &Config{Invisible: true})

	if user != admin.resources[0] {
		t.Error("resource not added")
	}

	if len(admin.GetMenus()) != 0 {
		t.Error("invisible resource registered in menu")
	}
}

func TestGetResource(t *testing.T) {
	admin := New(&qor.Config{DB: db})
	user := admin.AddResource(&User{})

	if admin.GetResource("User") != user {
		t.Error("resource not returned")
	}
}

func TestNewResource(t *testing.T) {
	admin := New(&qor.Config{DB: db})
	user := admin.NewResource(&User{})

	if user.Name != "User" {
		t.Error("default resource name didn't set")
	}
}

type UserWithCustomizedName struct{}

func (u *UserWithCustomizedName) ResourceName() string {
	return "CustomizedName"
}

func TestNewResourceWithCustomizedName(t *testing.T) {
	admin := New(&qor.Config{DB: db})
	user := admin.NewResource(&UserWithCustomizedName{})

	if user.Name != "CustomizedName" {
		t.Error("customize resource name didn't set")
	}
}
