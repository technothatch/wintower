//go:build !windows

package procs

func List() ([]Proc, error) {
	return nil, ErrUnsupported
}

func Kill(pid uint32) error {
	return ErrUnsupported
}
