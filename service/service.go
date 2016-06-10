package service

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/hoop33/entrevista"
)

// Service represents a service
type Service interface {
	Login() (string, error)
}

var services = make(map[string]Service)

func registerService(service Service) {
	parts := strings.Split(reflect.TypeOf(service).String(), ".")
	services[strings.ToLower(parts[len(parts)-1])] = service
}

// ForName returns the service for a given name
func ForName(name string) Service {
	if service, ok := services[strings.ToLower(name)]; ok {
		return service
	}
	return services["github"]
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
