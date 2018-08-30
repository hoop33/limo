package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hoop33/entrevista"
	"github.com/hoop33/limo/model"
	"github.com/xanzy/go-gitlab"
)

// Gitlab represents the Gitlab service
type Gitlab struct {
	insecure bool
}

// Login logs in to Gitlab
func (g *Gitlab) Login(ctx context.Context) (string, error) {
	interview := createInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitLab API token",
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

// AddStar stars a repo
func (g *Gitlab) AddStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	client := g.getClient(token)

	// Add the star
	project, _, err := client.Projects.StarProject(fmt.Sprintf("%s/%s", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	return model.NewStarFromGitlab(*project)
}

// DeleteStar unstars a repo
func (g *Gitlab) DeleteStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	client := g.getClient(token)

	// Add the star
	project, _, err := client.Projects.UnstarProject(fmt.Sprintf("%s/%s", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	return model.NewStarFromGitlab(*project)
}

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Gitlab) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
	client := g.getClient(token)

	currentPage := 1
	lastPage := 1

	starred := true
	for currentPage <= lastPage {
		projects, response, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Starred: &starred,
			ListOptions: gitlab.ListOptions{
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
			lastPage = response.TotalPages

			// Create a Star for each repository and put it on the channel
			for _, project := range projects {
				star, err := model.NewStarFromGitlab(*project)
				starChan <- &model.StarResult{
					Error: err,
					Star:  star,
				}
			}
		}
		// Go to the next page
		currentPage++
	}

	close(starChan)
}

// GetEvents returns the events for the authenticated user
func (g *Gitlab) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
	eventChan <- &model.EventResult{
		Error: errors.New("GitLab not yet supported"),
		Event: nil,
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Gitlab) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	close(trendingChan)
}

// SetInsecure sets whether to skip cert verification
func (g *Gitlab) SetInsecure(insecure bool) {
	g.insecure = insecure
}

func (g *Gitlab) getClient(token string) *gitlab.Client {
	client := &http.Client{
		// TODO make configurable
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: g.insecure,
			},
		},
	}
	return gitlab.NewClient(client, token)
}

func init() {
	registerService(&Gitlab{})
}
