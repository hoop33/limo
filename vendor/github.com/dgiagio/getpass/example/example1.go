package main

import (
	"fmt"
	"github.com/dgiagio/getpass"
)

func main() {
	pass, _ := getpass.GetPassword("Password: ")
	fmt.Printf("Entered password: %s\n", pass)
}
