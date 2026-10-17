//go:build windows

package services

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func statusString(s svc.State) string {
	switch s {
	case svc.Stopped:
		return "Stopped"
	case svc.StartPending:
		return "Starting"
	case svc.StopPending:
		return "Stopping"
	case svc.Running:
		return "Running"
	case svc.ContinuePending:
		return "ContinuePending"
	case svc.PausePending:
		return "PausePending"
	case svc.Paused:
		return "Paused"
	default:
		return "Unknown"
	}
}

const (
	startAutomatic = 0x2
	startManual    = 0x3
	startDisabled  = 0x4
)

func startTypeString(t uint32) string {
	switch t {
	case startAutomatic:
		return "Automatic"
	case startManual:
		return "Manual"
	case startDisabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}

func List() ([]Service, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("connect to service manager: %w", err)
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}

	out := make([]Service, 0, len(names))
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			continue
		}
		st, err := s.Query()
		if err != nil {
			s.Close()
			continue
		}
		cfg, err := s.Config()
		s.Close()
		if err != nil {
			cfg.DisplayName = ""
			cfg.StartType = 0xffffffff
		}
		out = append(out, Service{
			Name:        name,
			DisplayName: cfg.DisplayName,
			Status:      statusString(st.State),
			StartType:   startTypeString(cfg.StartType),
		})
	}
	return out, nil
}

func waitForState(s *mgr.Service, want svc.State) error {
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		st, err := s.Query()
		if err != nil {
			return err
		}
		if st.State == want {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for state %v", want)
}

func Control(name, action string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return err
	}
	defer s.Close()

	switch action {
	case "start":
		if err := s.Start(); err != nil {
			return err
		}
		return waitForState(s, svc.Running)
	case "stop":
		if _, err := s.Control(svc.Stop); err != nil {
			return err
		}
		return waitForState(s, svc.Stopped)
	case "restart":
		if _, err := s.Control(svc.Stop); err != nil {
			return err
		}
		if err := waitForState(s, svc.Stopped); err != nil {
			return err
		}
		if err := s.Start(); err != nil {
			return err
		}
		return waitForState(s, svc.Running)
	default:
		return fmt.Errorf("unknown action %q", action)
	}
}
