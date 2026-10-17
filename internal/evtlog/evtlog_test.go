package evtlog

import (
	"testing"
	"time"
	"unicode/utf8"
)

func TestParseJSONArray(t *testing.T) {
	data := []byte(`[
		{"Time":"2026-10-16T12:00:00.0000000Z","Id":7036,"LevelDisplayName":"Information","ProviderName":"Service Control Manager","Message":"The Windows Update service entered the running state."},
		{"Time":"2026-10-16T12:01:00.0000000Z","Id":41,"LevelDisplayName":"Critical","ProviderName":"Kernel-Power","Message":"The system has rebooted without cleanly shutting down first."}
	]`)
	evs, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON: %v", err)
	}
	if len(evs) != 2 {
		t.Fatalf("got %d events, want 2", len(evs))
	}
	if evs[0].ID != 7036 || evs[0].Provider != "Service Control Manager" {
		t.Errorf("unexpected first event: %+v", evs[0])
	}
	if !evs[1].Time.Equal(time.Date(2026, 10, 16, 12, 1, 0, 0, time.UTC)) {
		t.Errorf("unexpected time: %v", evs[1].Time)
	}
	if evs[1].Level != "Critical" {
		t.Errorf("unexpected level: %q", evs[1].Level)
	}
}

func TestParseJSONSingleObject(t *testing.T) {
	data := []byte(`{"Time":"2026-10-16T12:00:00Z","Id":1,"LevelDisplayName":"Error","ProviderName":"Test","Message":"boom"}`)
	evs, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON: %v", err)
	}
	if len(evs) != 1 || evs[0].Message != "boom" {
		t.Fatalf("unexpected result: %+v", evs)
	}
}

func TestTruncate(t *testing.T) {
	long := make([]rune, 500)
	for i := range long {
		long[i] = 'a'
	}
	evs, err := ParseJSON([]byte(`{"Id":1,"Message":"` + string(long) + `"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := utf8.RuneCountInString(evs[0].Message); got > 200 {
		t.Errorf("message not truncated: %d runes", got)
	}
}
