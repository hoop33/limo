package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/hoop33/entrevista"
	"github.com/hoop33/limo/model"
)

var errNotLoggedIn = errors.New("not logged in")

// Service represents a service
type Service interface {
	Login(ctx context.Context) (string, error)
	AddStar(ctx context.Context, token, owner, repo string) (*model.Star, error)
	DeleteStar(ctx context.Context, token, owner, repo string) (*model.Star, error)
	GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string)
	GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token, language string, verbose bool)
	GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int)
	SetInsecure(insecure bool)
}

var services = make(map[string]Service)

func registerService(service Service) {
	services[Name(service)] = service
}

// Name returns the name of a service
func Name(service Service) string {
	parts := strings.Split(reflect.TypeOf(service).String(), ".")
	return strings.ToLower(parts[len(parts)-1])
}

// ForName returns the service for a given name, or an error if it doesn't exist
func ForName(name string, insecure bool) (Service, error) {
	if service, ok := services[strings.ToLower(name)]; ok {
		service.SetInsecure(insecure)
		return service, nil
	}
	return &NotFound{}, fmt.Errorf("service '%s' not found", name)
}

func createInterview() *entrevista.Interview {
	interview := entrevista.NewInterview()
	interview.ShowOutput = func(message string) {
		fmt.Print(color.GreenString(message))
	}
	interview.ShowError = func(message string) {
		color.Red(message)
	}
	return interview
}
