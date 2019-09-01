package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/hoop33/entrevista"
	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/segmentio/stats/httpstats"
	"golang.org/x/oauth2"
)

// Github represents the Github service
type Github struct {
	c httpcache.Cache
}

// Login logs in to Github
func (g *Github) Login(ctx context.Context) (string, error) {
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

// AddStar stars a repo
func (g *Github) AddStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	client := g.getClient(token)

	// Add the star
	_, err := client.Activity.Star(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	// Get the repo details
	r, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return model.NewStarFromGithub(&github.Timestamp{Time: time.Now()}, *r)
}

// DeleteStar unstars a repo
func (g *Github) DeleteStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	client := g.getClient(token)

	// Remove the star
	_, err := client.Activity.Unstar(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	// Get the repo details
	r, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return model.NewStarFromGithub(&github.Timestamp{Time: time.Now()}, *r)
}

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string) {
	defer close(starChan)

	client := g.getClient(token)

	// The first response will give us the correct value for the last page
	currentPage := 1
	lastPage := 1

	for currentPage <= lastPage {
		repos, response, err := client.Activity.ListStarred(ctx, user, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    currentPage,
				PerPage: 100,
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
				star, err := model.NewStarFromGithub(repo.StarredAt, *repo.Repository)
				starChan <- &model.StarResult{
					Error: err,
					Star:  star,
				}
			}
		}
		// Go to the next page
		currentPage++
	}
}

// GetEvents returns the events for the authenticated user
func (g *Github) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
	client := g.getClient(token)

	currentPage := page
	lastPage := page + count - 1

	for currentPage <= lastPage {
		events, _, err := client.Activity.ListEventsReceivedByUser(ctx, user, false, &github.ListOptions{
			Page:    currentPage,
			PerPage: 100,
		})

		if err != nil {
			eventChan <- &model.EventResult{
				Error: err,
				Event: nil,
			}
		} else {
			for _, event := range events {
				eventChan <- &model.EventResult{
					Error: nil,
					Event: model.NewEventFromGithub(event),
				}
			}
		}
		currentPage++
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Github) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	client := g.getClient(token)

	// TODO perhaps allow them to specify multiple pages?
	// Might be overkill -- first page probably plenty

	// TODO Make this more configurable. Sort by stars, forks, default.
	// Search by number of stars, pushed, created, or whatever.
	// Lots of possibilities.

	q := g.getDateSearchString()

	if language != "" {
		q = fmt.Sprintf("language:%s %s", language, q)
	}

	if verbose {
		fmt.Println("q =", q)
	}

	result, _, err := client.Search.Repositories(ctx, q, &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	})

	// If we got an error, put it on the channel
	if err != nil {
		trendingChan <- &model.StarResult{
			Error: err,
			Star:  nil,
		}
	} else {
		// Create a Star for each repository and put it on the channel
		for _, repo := range result.Repositories {
			star, err := model.NewStarFromGithub(nil, repo)
			trendingChan <- &model.StarResult{
				Error: err,
				Star:  star,
			}
		}
	}

	close(trendingChan)
}

// SetInsecure sets whether to skip cert verification
func (g *Github) SetInsecure(insecure bool) {
}

func (g *Github) getDateSearchString() string {
	// TODO make this configurable
	// Default should be in configuration file
	// and should be able to override from command line
	// TODO should be able to specify whether "created" or "pushed"
	date := time.Now().Add(-7 * (24 * time.Hour))
	return fmt.Sprintf("created:>%s", date.Format("2006-01-02"))
}

func (g *Github) getClient(token string) *github.Client {

	// var err error
	if g.c == nil {
		g.c = diskcache.New(config.ConfigDirectoryPath + "/cache/http")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	var hc http.Client

	t := httpcache.NewTransport(g.c)
	t.MarkCachedResponses = true

	hc.Transport = httpstats.NewTransport(t)
	timeout := time.Duration(10 * time.Second)

	return github.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Base:   hc.Transport,
			Source: ts,
		},
		Timeout: timeout,
	})
}

func (g *Github) getClientOld(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func init() {
	registerService(&Github{})
}
