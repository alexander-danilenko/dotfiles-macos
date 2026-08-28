package main

import (
	"testing"
	"time"
)

func TestPaceColor(t *testing.T) {
	const window = 5 * 60 * 60
	for _, c := range []struct {
		name           string
		usedPct, gone  float64 // gone is the share of the window that has passed
		want, wantName string
	}{
		{"on pace", 10, 0.10, green, "green"},
		{"under pace", 25, 0.50, green, "green"},
		{"burst early", 50, 0.10, red, "red"},
		{"the same burst, an hour of idling later", 50, 0.30, orange, "orange"},
		{"and later still", 50, 0.45, yello, "yellow"},
		{"nearly spent but nearly reset", 90, 0.95, green, "green"},
		{"a trickle in the first minute", 2, 0.003, green, "green"},
	} {
		if got := paceColor(c.usedPct, window*(1-c.gone), window); got != c.want {
			t.Errorf("%s: %.0f%% used with %.0f%% of the window gone is %q, want %s",
				c.name, c.usedPct, c.gone*100, got, c.wantName)
		}
	}
	// A 7 day window takes the same treatment, only the length differs.
	if got := paceColor(50, (7*24*time.Hour).Seconds()*0.9, (7 * 24 * time.Hour).Seconds()); got != red {
		t.Errorf("half the weekly limit in a tenth of the week is %q, want red", got)
	}
}

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
