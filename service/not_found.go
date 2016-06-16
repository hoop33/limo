package service

import "errors"

type NotFound struct {
}

func (nf *NotFound) Login() (string, error) {
	return "", errors.New("Service not found")
}
