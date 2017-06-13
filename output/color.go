package output

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
)

const defaultInterval = 300
const minInterval = 250
const defaultColor = "yellow"

var spin *spinner.Spinner
var cfg *config.OutputConfig

// Color is a color text output
type Color struct {
}

// Configure configures the output
func (c *Color) Configure(oc *config.OutputConfig) {
	cfg = oc
}

// Inline displays text in line
func (c *Color) Inline(s string) {
	fmt.Print(color.GreenString(s))
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

// Event displays an event {
func (c *Color) Event(event *model.Event) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.YellowString(event.Who))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.GreenString(fmt.Sprintf(" %s", event.What)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.BlueString(fmt.Sprintf(" %s", event.Which)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.RedString(fmt.Sprintf(" (%s)", event.URL)))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.MagentaString(fmt.Sprintf(" %s", humanize.Time(event.When))))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// StarLine displays a star in one line
func (c *Color) StarLine(star *model.Star) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(*star.FullName))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", star.Stargazers)))
	if err != nil {
		c.Error(err.Error())
	}

	if star.Language != nil {
		_, err = buffer.WriteString(color.GreenString(fmt.Sprintf(" %s", *star.Language)))
		if err != nil {
			c.Error(err.Error())
		}
	}

	if star.URL != nil {
		_, err = buffer.WriteString(color.RedString(fmt.Sprintf(" %s", *star.URL)))
		if err != nil {
			c.Error(err.Error())
		}
	}

	fmt.Println(buffer.String())
}

// Star displays a star
func (c *Color) Star(star *model.Star) {
	c.StarLine(star)

	if len(star.Tags) > 0 {
		var buffer bytes.Buffer
		leader := ""
		for _, tag := range star.Tags {
			_, err := buffer.WriteString(color.MagentaString(fmt.Sprintf("%s%s", leader, tag.Name)))
			if err != nil {
				c.Error(err.Error())
			}
			leader = ", "
		}
		fmt.Println(buffer.String())
	}

	if star.Description != nil && *star.Description != "" {
		color.White(*star.Description)
	}

	if star.Homepage != nil && *star.Homepage != "" {
		color.Red(fmt.Sprintf("Home page: %s", *star.Homepage))
	}

	color.Green(fmt.Sprintf("Starred on %s", star.StarredAt.Format(time.UnixDate)))
}

// Tag displays a tag
func (c *Color) Tag(tag *model.Tag) {
	var buffer bytes.Buffer

	_, err := buffer.WriteString(color.BlueString(tag.Name))
	if err != nil {
		c.Error(err.Error())
	}

	_, err = buffer.WriteString(color.YellowString(fmt.Sprintf(" ★ :%d", tag.StarCount)))
	if err != nil {
		c.Error(err.Error())
	}

	fmt.Println(buffer.String())
}

// Tick displays evidence that the program is working
func (c *Color) Tick() {
	if spin == nil {
		index := 0
		interval := defaultInterval
		clr := defaultColor
		if cfg != nil {
			index = cfg.SpinnerIndex
			if index < 0 || index > len(spinner.CharSets) {
				index = 0
			}
			interval = cfg.SpinnerInterval
			if interval < minInterval {
				interval = minInterval
			}
			clr = cfg.SpinnerColor
			if clr == "" {
				clr = defaultColor
			}
		}
		spin = spinner.New(spinner.CharSets[index], time.Duration(interval)*time.Millisecond)
		spin.Suffix = color.CyanString(" Updating")
		if err := spin.Color(clr); err != nil {
			c.Error(err.Error())
		}
	}
}

func init() {
	registerOutput(&Color{})
}
