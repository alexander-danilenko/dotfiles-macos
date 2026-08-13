package main

import "testing"

func TestModelName(t *testing.T) {
	for _, c := range []struct{ id, family, version string }{
		{"claude-opus-5", "opus", "5"},
		{"claude-opus-5[1m]", "opus", "5"},
		{"claude-opus-4-7", "opus", "4.7"},
		{"claude-3-5-sonnet-20241022", "sonnet", "3.5"},
		{"claude-haiku-4-5-20251001", "haiku", "4.5"},
		{"claude-fable-5", "fable", "5"},
		{"gpt-9", "unknown", ""},
	} {
		family, _ := modelFamily(c.id)
		version := modelVersion(c.id, family)
		if family != c.family || version != c.version {
			t.Errorf("%s: got %s/%s, want %s/%s", c.id, family, version, c.family, c.version)
		}
	}
}
