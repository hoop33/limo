package service

import (
	"context"
	"errors"

	"github.com/hoop33/limo/model"
)

var errNotFound = errors.New("service not found")

// NotFound is used when the specified service is not found
type NotFound struct {
}

// Login is not implemented
func (nf *NotFound) Login(ctx context.Context) (string, error) {
	return "", errNotFound
}

// AddStar is not implemented
func (nf *NotFound) AddStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	return nil, errNotFound
}

// DeleteStar is not implemented
func (nf *NotFound) DeleteStar(ctx context.Context, token, owner, repo string) (*model.Star, error) {
	return nil, errNotFound
}

// GetStars is not implemented
func (nf *NotFound) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
}

// GetTrending is not implemented
func (nf *NotFound) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
}

// GetEvents is not implemented
func (nf *NotFound) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
}

// SetInsecure sets the service to skip cert verification
func (nf *NotFound) SetInsecure(insecure bool) {
}
