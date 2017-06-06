package model

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

const defaultWho = "somebody"
const defaultWhat = "did something with"
const defaultWhich = "some repository"

var eventTypes = map[string]string{
	"CreateEvent":       "created",
	"DeleteEvent":       "deleted",
	"ForkEvent":         "forked",
	"IssueCommentEvent": "commented on an issue on",
	"IssuesEvent":       "opened an issue on",
	"PullRequestEvent":  "opened a pull request on",
	"PushEvent":         "pushed to",
	"WatchEvent":        "starred",
}

// Event is a git-hosting service event
type Event struct {
	Who   string
	What  string
	Which string
	URL   string
	When  time.Time
}

// EventResult wraps an event and an error
type EventResult struct {
	Event *Event
	Error error
}

// NewEventFromGithub creates an Event from a Github event
func NewEventFromGithub(event *github.Event) *Event {
	who := defaultWho
	if event.Actor != nil && event.Actor.Login != nil {
		who = *event.Actor.Login
	}

	what := defaultWhat
	if event.Type != nil {
		if action, ok := eventTypes[*event.Type]; ok {
			what = action
		}
	}

	which := defaultWhich
	url := ""
	if event.Repo != nil {
		if event.Repo.Name != nil {
			which = *event.Repo.Name
			url = fmt.Sprintf("https://github.com/%s", which)
		}
	}

	when := time.Now()
	if event.CreatedAt != nil {
		when = *event.CreatedAt
	}

	return &Event{
		Who:   who,
		What:  what,
		Which: which,
		URL:   url,
		When:  when,
	}
}
