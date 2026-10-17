//go:build !windows

package services

func List() ([]Service, error) {
	return nil, ErrUnsupported
}

func Control(name, action string) error {
	return ErrUnsupported
}
