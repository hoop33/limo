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

	if star.Stargazers > 0 {
		_, err = buffer.WriteString(fmt.Sprintf(" (*: %d)", star.Stargazers))
		if err != nil {
			t.Error(err.Error())
		}
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
	t.StarLine(star)

	if len(star.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range star.Tags {
			_, err := buffer.WriteString(fmt.Sprintf("%s%s", leader, tag.Name))
			if err != nil {
				t.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if star.Description != nil && *star.Description != "" {
		fmt.Println(*star.Description)
	}

	if star.Homepage != nil && *star.Homepage != "" {
		fmt.Println(fmt.Sprintf("Home page: %s", *star.Homepage))
	}

	if star.URL != nil {
		fmt.Println(fmt.Sprintf("URL: %s", *star.URL))
	}
}

// Tick displays evidence that the program is working
func (t *Text) Tick() {
	fmt.Print(".")
}

func init() {
	registerOutput(&Text{})
}
