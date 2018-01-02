package service

import (
	"context"
	"errors"

	"github.com/hoop33/limo/model"
)

// NotFound is used when the specified service is not found
type NotFound struct {
}

// Login is not implemented
func (nf *NotFound) Login(ctx context.Context) (string, error) {
	return "", errors.New("service not found")
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
