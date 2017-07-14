package entrevista_test

import (
	"fmt"

	"github.com/hoop33/entrevista"
)

func Example() {
	interview := entrevista.NewInterview()
	interview.ReadAnswer = func(question *entrevista.Question) (string, error) {
		return question.Key, nil
	}
	interview.Questions = []entrevista.Question{
		{
			Key:      "name",
			Text:     "Enter your name",
			Required: true,
		},
		{
			Key:           "email",
			Text:          "Enter your email address",
			DefaultAnswer: "john.doe@example.com",
		},
	}
	answers, err := interview.Run()

	if err == nil {
		fmt.Print(answers["name"], ",", answers["email"])
	} else {
		fmt.Print(err.Error())
	}
	// Output: Enter your name: Enter your email address (john.doe@example.com): name,email
}
