//go:build !windows

package evtlog

func Events() ([]Event, error) {
	return nil, ErrUnsupported
}
