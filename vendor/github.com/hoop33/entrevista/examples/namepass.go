package main

import (
	"fmt"
	"log"

	"github.com/hoop33/entrevista"
)

func main() {
	interview := entrevista.NewInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "name",
			Text:     "Enter your name",
			Required: true,
		},
		{
			Key:      "password",
			Text:     "Enter your password",
			Hidden:   true,
			Required: true,
		},
		*entrevista.NewBoolQuestion("show", "Should I display the password?"),
		*entrevista.NewNumberQuestion("times", "How many times should I show your name?"),
	}

	answers, err := interview.Run()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < answers["times"].(int); i++ {
		fmt.Println(answers["name"])
	}

	if answers["show"].(bool) {
		fmt.Println(answers["password"])
	}
}
