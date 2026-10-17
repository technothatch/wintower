//go:build windows

package procs

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	psapi                = windows.NewLazySystemDLL("psapi.dll")
	procGetProcessMemory = psapi.NewProc("GetProcessMemoryInfo")
)

type processMemoryCountersEx struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr
}

func processName(pid uint32) string {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)

	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return ""
	}
	return filepath.Base(windows.UTF16ToString(buf[:size]))
}

func processMem(h windows.Handle) uint64 {
	var counters processMemoryCountersEx
	counters.CB = uint32(unsafe.Sizeof(counters))
	r, _, _ := procGetProcessMemory.Call(uintptr(h), uintptr(unsafe.Pointer(&counters)), unsafe.Sizeof(counters))
	if r == 0 {
		return 0
	}
	return uint64(counters.WorkingSetSize)
}

func List() ([]Proc, error) {
	buf := make([]uint32, 4096)
	for {
		var read uint32
		if err := windows.EnumProcesses(buf, &read); err != nil {
			return nil, fmt.Errorf("enum processes: %w", err)
		}
		if int(read) < len(buf)*4 {
			buf = buf[:int(read)/4]
			break
		}
		buf = make([]uint32, len(buf)*2)
	}

	out := make([]Proc, 0, len(buf))
	for _, pid := range buf {
		if pid == 0 {
			continue
		}
		p := Proc{PID: pid, Name: processName(pid)}
		if p.Name == "" {
			continue
		}
		if h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid); err == nil {
			p.MemBytes = processMem(h)
			windows.CloseHandle(h)
		}
		out = append(out, p)
	}
	return out, nil
}

func Kill(pid uint32) error {
	h, err := windows.OpenProcess(windows.PROCESS_TERMINATE|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return fmt.Errorf("open process %d: %w", pid, err)
	}
	defer windows.CloseHandle(h)
	if err := windows.TerminateProcess(h, 1); err != nil {
		return fmt.Errorf("terminate process %d: %w", pid, err)
	}
	return nil
}
