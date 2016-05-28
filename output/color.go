package output

import "github.com/fatih/color"

// Color is a color text output
type Color struct {
}

// Info displays information
func (c Color) Info(s string) {
	color.Green(s)
}

// Error displays an error
func (c Color) Error(s string) {
	color.Red(s)
}

func init() {
	registerOutput(Color{})
}
