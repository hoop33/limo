package output

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hoop33/limo/model"
)

// Text is a monochrome text output
type Text struct {
}

// Info displays information
func (t *Text) Info(s string) {
	fmt.Println(s)
}

// Error displays an error
func (t *Text) Error(s string) {
	fmt.Fprintln(os.Stderr, s)
}

// Fatal displays an error and ends the program
func (t *Text) Fatal(s string) {
	t.Error(s)
	os.Exit(1)
}

// StarLine displays a star in one line
func (t *Text) StarLine(star *model.Star) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(*star.FullName)
	if err != nil {
		t.Error(err.Error())
	}
	if star.Language != nil {
		_, err := buffer.WriteString(fmt.Sprintf(" (%s)", *star.Language))
		if err != nil {
			t.Error(err.Error())
		}
	}
	fmt.Println(buffer.String())
}

// Star displays a star
func (t *Text) Star(star *model.Star) {
}

// Tick displays evidence that the program is working
func (t *Text) Tick() {
	fmt.Print(".")
}

func init() {
	registerOutput(&Text{})
}
