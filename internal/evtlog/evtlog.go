package evtlog

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrUnsupported = errors.New("wintower: event log is only available on Windows")

type Event struct {
	Time     time.Time
	ID       int
	Level    string
	Provider string
	Message  string
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

func ParseJSON(data []byte) ([]Event, error) {
	type rawEvent struct {
		TimeCreated      string `json:"Time"`
		ID               int    `json:"Id"`
		LevelDisplayName string `json:"LevelDisplayName"`
		ProviderName     string `json:"ProviderName"`
		Message          string `json:"Message"`
	}

	var raws []rawEvent
	if err := json.Unmarshal(data, &raws); err != nil {
		var one rawEvent
		if err2 := json.Unmarshal(data, &one); err2 != nil {
			return nil, err
		}
		raws = []rawEvent{one}
	}

	out := make([]Event, 0, len(raws))
	for _, r := range raws {
		ev := Event{
			ID:       r.ID,
			Level:    r.LevelDisplayName,
			Provider: r.ProviderName,
			Message:  truncate(r.Message, 200),
		}
		if t, err := time.Parse(time.RFC3339, r.TimeCreated); err == nil {
			ev.Time = t
		}
		out = append(out, ev)
	}
	return out, nil
}
