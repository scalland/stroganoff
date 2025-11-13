package version

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Fatal("GetVersion returned empty string")
	}
}

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Fatal("Version should not be empty")
	}

	if info.String() == "" {
		t.Fatal("String() should not be empty")
	}
}
