package output

import (
	"fmt"
	"os"
)

// Text is a monochrome text output
type Text struct {
}

// Info displays information
func (t Text) Info(s string) {
	fmt.Println(s)
}

// Error displays an error
func (t Text) Error(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func init() {
	registerOutput(Text{})
}
