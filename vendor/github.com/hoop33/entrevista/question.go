package entrevista

import (
	"fmt"
	"reflect"
	"regexp"
)

// Question is a question in an interview
type Question struct {
	// The key for the answer map for the answer. Required.
	Key string
	// The text of the question. Required.
	Text string
	// The type of the expected answer.
	AnswerKind reflect.Kind
	// Whether an answer is required.
	Required bool
	// The default answer.
	DefaultAnswer string
	// The regular expression used to validate the answer.
	RegularExpression *regexp.Regexp
	// Whether to validate minimum and maximum
	ValidateMinMax bool
	// The minimum (length for a string, value for a number)
	Minimum int
	// The maximum (length for a string, value for a number)
	Maximum int
	// The error message to display if the answer is required and not supplied.
	RequiredMessage string
	// The error message to display if the answer is invalid.
	InvalidMessage string
	// Whether the response should be hidden when typed
	Hidden bool
}

// YesNoRegularExpression is a regular expression that accepts Y, y, N, or n
var YesNoRegularExpression = regexp.MustCompile("^[YyNn]")

// NewQuestion creates a new question
func NewQuestion(key string, text string) *Question {
	return &Question{
		Key:        key,
		Text:       text,
		AnswerKind: reflect.String,
	}
}

// NewStringQuestion creates a new question that expects a string response
func NewStringQuestion(key string, text string, minLength int, maxLength int) *Question {
	return &Question{
		Key:            key,
		Text:           text,
		AnswerKind:     reflect.String,
		ValidateMinMax: true,
		Minimum:        minLength,
		Maximum:        maxLength,
		InvalidMessage: fmt.Sprintf("Your answer must have a length between %d and %d characters.", minLength, maxLength),
	}
}

// NewBoolQuestion creates a new question that expects a boolean response
func NewBoolQuestion(key string, text string) *Question {
	return &Question{
		Key:               key,
		Text:              text,
		AnswerKind:        reflect.Bool,
		DefaultAnswer:     "no",
		RegularExpression: YesNoRegularExpression,
		InvalidMessage:    "You must answer yes or no.",
	}
}

// NewNumberQuestion creates a new question that expects a number response
func NewNumberQuestion(key string, text string) *Question {
	return &Question{
		Key:            key,
		Text:           text,
		AnswerKind:     reflect.Int,
		ValidateMinMax: false,
	}
}

// NewNumberInRangeQuestion creates a new questions that expects a number response in a range
func NewNumberInRangeQuestion(key string, text string, min int, max int) *Question {
	return &Question{
		Text:           text,
		AnswerKind:     reflect.Int,
		ValidateMinMax: true,
		Minimum:        min,
		Maximum:        max,
		InvalidMessage: fmt.Sprintf("You must enter a number between %d and %d.", min, max),
	}
}
