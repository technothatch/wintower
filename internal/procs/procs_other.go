//go:build !windows

package procs

func List() ([]Proc, error) {
	return nil, ErrUnsupported
}

func Kill(pid uint32) error {
	return ErrUnsupported
}

func SetPriority(pid uint32, class uint32) error {
	return ErrUnsupported
}
