package service

import (
	"errors"

	"github.com/hoop33/limo/model"
)

// NotFound is used when the specified service is not found
type NotFound struct {
}

// Login is not implemented
func (nf *NotFound) Login() (string, error) {
	return "", errors.New("Service not found")
}

// GetStars is not implemented
func (nf *NotFound) GetStars(starChan chan<- *model.StarResult, token string, user string) {
}
