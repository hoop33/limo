package entrevista_test

import (
	"reflect"
	"testing"

	"github.com/hoop33/entrevista"
	"github.com/stretchr/testify/assert"
)

func MirrorAnswer(question *entrevista.Question) (string, error) {
	return question.Text, nil
}

func QuietAnswer(question *entrevista.Question) (string, error) {
	return "", nil
}

func QuietOutput(message string) {
}

func QuietError(message string) {
}

func TestAnswersGetFilled(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = MirrorAnswer

	interview.Questions = []entrevista.Question{
		{
			Key:  "one",
			Text: "One",
		},
		{
			Key:  "two",
			Text: "Two",
		},
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["one"], "One")
	assert.Equal(t, answers["two"], "Two")
}

func TestDefaultAnswerIsReturnedForBlank(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = QuietAnswer

	interview.Questions = []entrevista.Question{
		{
			Key:           "one",
			Text:          "One",
			DefaultAnswer: "First default",
		},
		{
			Key:           "two",
			Text:          "Two",
			DefaultAnswer: "Second default",
		},
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["one"], "First default")
	assert.Equal(t, answers["two"], "Second default")
}

func TestBooleansAreReturnedForBooleans(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = MirrorAnswer

	interview.Questions = []entrevista.Question{
		{
			Key:        "yes",
			Text:       "Yes",
			AnswerKind: reflect.Bool,
		},
		{
			Key:        "no",
			Text:       "No",
			AnswerKind: reflect.Bool,
		},
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["yes"], true)
	assert.Equal(t, answers["no"], false)
}

func TestBooleanAnswersActAsBooleans(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = MirrorAnswer

	interview.Questions = []entrevista.Question{
		*entrevista.NewBoolQuestion("yes", "yes"),
		*entrevista.NewBoolQuestion("no", "no"),
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["yes"], true)
	assert.Equal(t, answers["no"], false)
}

func TestNumbersAreReturnedForNumbers(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = MirrorAnswer

	interview.Questions = []entrevista.Question{
		{
			Key:        "1",
			Text:       "12345",
			AnswerKind: reflect.Int,
		},
		{
			Key:        "2",
			Text:       "43",
			AnswerKind: reflect.Int,
		},
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["1"], 12345)
	assert.Equal(t, answers["2"], 43)
}

func TestMinAndMaxWorkForNumbers(t *testing.T) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = QuietOutput
	interview.ShowError = QuietError
	interview.ReadAnswer = MirrorAnswer

	interview.Questions = []entrevista.Question{
		{
			Key:        "max",
			Text:       "12345",
			AnswerKind: reflect.Int,
			Maximum:    99999,
		},
		{
			Key:        "min",
			Text:       "-43",
			AnswerKind: reflect.Int,
			Minimum:    -50,
		},
	}
	answers, err := interview.Run()
	assert.Equal(t, err, nil)
	assert.Equal(t, answers["max"], 12345)
	assert.Equal(t, answers["min"], -43)
}
