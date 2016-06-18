package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Color is a color text output
type Color struct {
}

// Info displays information
func (c *Color) Info(s string) {
	color.Green(s)
}

// Error displays an error
func (c *Color) Error(s string) {
	color.Red(s)
}

// Fatal displays an error and ends the program
func (c *Color) Fatal(s string) {
	c.Error(s)
	os.Exit(1)
}

// Tick displays evidence that the program is working
func (c *Color) Tick() {
	fmt.Print(color.CyanString("."))
}

func init() {
	registerOutput(&Color{})
}
