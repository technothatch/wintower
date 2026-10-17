package services

import "errors"

var ErrUnsupported = errors.New("wintower: services are only available on Windows")

type Service struct {
	Name        string
	DisplayName string
	Status      string
	StartType   string
}
