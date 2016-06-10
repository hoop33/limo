package service

import "github.com/hoop33/entrevista"

type Github struct {
}

func (g *Github) Login() (string, error) {
	interview := createInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitHub API token",
			Required: true,
			Hidden:   true,
		},
	}

	answers, err := interview.Run()
	if err != nil {
		return "", err
	}
	return answers["token"].(string), nil
}

func init() {
	registerService(&Github{})
}
