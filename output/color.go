package output

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hoop33/limo/model"
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

// StarLine displays a star in one line
func (c *Color) StarLine(star *model.Star) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(*star.FullName))
	if err != nil {
		c.Error(err.Error())
	}

	if star.Language != nil {
		_, err := buffer.WriteString(color.YellowString(fmt.Sprintf(" (%s)", *star.Language)))
		if err != nil {
			c.Error(err.Error())
		}
	}
	fmt.Println(buffer.String())
}

// Star displays a star
func (c *Color) Star(star *model.Star) {
}

// Tick displays evidence that the program is working
func (c *Color) Tick() {
	fmt.Print(color.CyanString("."))
}

func init() {
	registerOutput(&Color{})
}
