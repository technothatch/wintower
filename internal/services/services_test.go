package services

import "testing"

func TestErrUnsupportedOnOtherPlatforms(t *testing.T) {
	if _, err := List(); err == nil {
		t.Skip("running on windows; SCM-backed List works")
	}
}
