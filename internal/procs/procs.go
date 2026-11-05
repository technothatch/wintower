package procs

import "errors"

var ErrUnsupported = errors.New("wintower: process listing is only available on Windows")

type Proc struct {
	PID      uint32
	Name     string
	MemBytes uint64
	Priority string
}

// Win32 priority class values (PROCESS_INFORMATION_CLASS / SetPriorityClass).
const (
	PriorityIdle        uint32 = 0x00000040
	PriorityBelowNormal uint32 = 0x00004000
	PriorityNormal      uint32 = 0x00000020
	PriorityAboveNormal uint32 = 0x00008000
	PriorityHigh        uint32 = 0x00000080
	PriorityRealtime    uint32 = 0x00000100
)

type PriorityOption struct {
	Name  string
	Class uint32
}

// Priorities lists the classes offered in the UI. Realtime is deliberately
// omitted: it can starve the whole machine if misused.
var Priorities = []PriorityOption{
	{Name: "Idle", Class: PriorityIdle},
	{Name: "Below Normal", Class: PriorityBelowNormal},
	{Name: "Normal", Class: PriorityNormal},
	{Name: "Above Normal", Class: PriorityAboveNormal},
	{Name: "High", Class: PriorityHigh},
}

func PriorityName(class uint32) string {
	switch class {
	case PriorityIdle:
		return "Idle"
	case PriorityBelowNormal:
		return "Below Normal"
	case PriorityNormal:
		return "Normal"
	case PriorityAboveNormal:
		return "Above Normal"
	case PriorityHigh:
		return "High"
	case PriorityRealtime:
		return "Realtime"
	default:
		return ""
	}
}
