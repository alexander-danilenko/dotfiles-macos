package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestPace(t *testing.T) {
	const window = 5 * 60 * 60
	for _, c := range []struct {
		name          string
		usedPct, gone float64 // gone is the share of the window that has passed
		want          float64
		wantColor     string
		colorName     string
	}{
		{"on pace", 10, 0.10, 1.0, green, "green"},
		{"under pace", 25, 0.50, 0.5, green, "green"},
		{"burst early", 50, 0.10, 5.0, red, "red"},
		{"the same burst, an hour of idling later", 50, 0.30, 1.667, orange, "orange"},
		{"and later still", 50, 0.45, 1.111, yello, "yellow"},
		{"nearly spent but nearly reset", 90, 0.95, 0.947, green, "green"},
		{"a trickle in the first minute", 2, 0.003, 0.4, green, "green"},
		{"the whole limit at once, the ceiling", 100, 0, 20, red, "red"},
		{"an untouched window", 0, 0, 0, green, "green"},
	} {
		got := pace(c.usedPct, window*(1-c.gone), window)
		if math.Abs(got-c.want) > 0.001 {
			t.Errorf("%s: %.0f%% used with %.0f%% of the window gone paces %.3f, want %.3f",
				c.name, c.usedPct, c.gone*100, got, c.want)
		}
		if color := paceColor(got); color != c.wantColor {
			t.Errorf("%s: pace %.3f is %q, want %s", c.name, got, color, c.colorName)
		}
	}
	// A 7 day window takes the same treatment, only the length differs.
	week := (7 * 24 * time.Hour).Seconds()
	if got := pace(50, week*0.9, week); math.Abs(got-5) > 0.001 {
		t.Errorf("half the weekly limit in a tenth of the week paces %.3f, want 5", got)
	}
	// A zero window would divide by zero, so it falls back to the plain usage.
	if got := pace(50, 0, 0); got != 0.5 {
		t.Errorf("a zero-length window paces %v, want 0.5", got)
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

func TestGaugeBar(t *testing.T) {
	// The leading edge slants only while the bar is still moving.
	for _, c := range []struct {
		usedPct float64
		want    string
	}{
		{0, "··········"},
		{50, "████◤·····"},
		{95, "████████◤·"},
		{100, "██████████"},
	} {
		if got := gauge(c.usedPct, "x", "", ""); !strings.Contains(got, c.want) {
			t.Errorf("%.0f%% renders %q, want a bar of %q", c.usedPct, got, c.want)
		}
	}
}
