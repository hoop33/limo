package service

import (
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/hoop33/entrevista"
	"github.com/hoop33/limo/model"
)

// Github represents the Github service
type Github struct {
}

// Login logs in to Github
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

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(starChan chan<- *model.StarResult, token string, user string) {
	client := getClient(token)

	// The first response will give us the correct value for the last page
	currentPage := 1
	lastPage := 1

	for currentPage <= lastPage {
		repos, response, err := client.Activity.ListStarred(user, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page: currentPage,
			},
		})
		// If we got an error, put it on the channel
		if err != nil {
			starChan <- &model.StarResult{
				Error: err,
				Star:  nil,
			}
		} else {
			// Set last page only if we didn't get an error
			lastPage = response.LastPage

			// Create a Star for each repository and put it on the channel
			for _, repo := range repos {
				starChan <- &model.StarResult{
					Error: nil,
					Star:  model.NewStarFromGithub(*repo.Repository),
				}
			}
		}
		// Go to the next page
		currentPage++
	}
	close(starChan)
}

func getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func init() {
	registerService(&Github{})
}
