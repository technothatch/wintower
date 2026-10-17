package procs

import "errors"

var ErrUnsupported = errors.New("wintower: process listing is only available on Windows")

type Proc struct {
	PID      uint32
	Name     string
	MemBytes uint64
}
