// Claude Code status line, run via `go run`. Parses the JSON on stdin once,
// shells out only for git, and renders through text/template.
// Template:
//
//	Line 1: {{pwd blue}} on {{git_branch green}} [+N|-N]
//	Line 2: {model}:{effort} | ctx:{used}/{total} [| $cost][| 5h:{pct}% bar][| 7d:{pct}% bar]
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// --- ANSI helpers ---
const (
	reset = "\x1b[0m"
	green = "\x1b[32m"
	yello = "\x1b[33m"
	red   = "\x1b[31m"
	blue  = "\x1b[34m"
	money = "\x1b[38;5;78m" // soft green for the API-equivalent $ cost

	// RPG rarity colors (256-color) for model families.
	// Rare → Haiku, Epic → Sonnet, Legendary → Opus.
	// Danger (crimson) → Fable, Mythos: high-cost models, flagged visually.
	rarityRare      = "\x1b[38;5;69m"        // Haiku
	rarityEpic      = "\x1b[38;5;135m"       // Sonnet
	rarityLegendary = "\x1b[38;5;220m"       // Opus
	colorDanger     = "\x1b[38;2;220;20;60m" // Fable / Mythos — crimson #DC143C (cost warning)
	colorFallback   = "\x1b[38;5;246m"       // unknown/fallback
)

type limit struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *int64   `json:"resets_at"`
}

type input struct {
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	Model struct {
		ID string `json:"id"`
	} `json:"model"`
	Effort struct {
		Level string `json:"level"`
	} `json:"effort"`
	ContextWindow struct {
		UsedPercentage    *float64 `json:"used_percentage"`
		ContextWindowSize *float64 `json:"context_window_size"`
	} `json:"context_window"`
	// Field name has varied across Claude Code versions; probe the known spellings.
	Cost struct {
		TotalCostUsd *float64 `json:"total_cost_usd"`
		TotalCostCC  *float64 `json:"totalCost"`
		TotalCost    *float64 `json:"total_cost"`
	} `json:"cost"`
	RateLimits struct {
		FiveHour limit `json:"five_hour"`
		SevenDay limit `json:"seven_day"`
	} `json:"rate_limits"`
}

// view is what the template renders. Every field is already colored.
type view struct {
	Dir    string
	Branch string
	Diff   string
	Model  string
	Ctx    string
	Cost   string
	Limits []string
}

const layout = `{{.Dir}}{{.Branch}} {{.Diff}}` +
	`{{if .Model}}
{{.Model}}{{with .Ctx}} | ctx:{{.}}{{end}}{{with .Cost}} | {{.}}{{end}}{{range .Limits}} | {{.}}{{end}}{{end}}`

// Run git quietly; return "" on any failure (not a repo, git missing, etc.).
func git(cwd string, args ...string) string {
	out, err := exec.Command("git", append([]string{"-C", cwd, "--no-optional-locks"}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Color for a percentage: <=60 green, <=80 yellow, else red.
func pctColor(pct int) string {
	switch {
	case pct <= 60:
		return green
	case pct <= 80:
		return yello
	default:
		return red
	}
}

var (
	numericPrefix = regexp.MustCompile(`^[0-9]+(?:-[0-9]+)*`)
	majorMinor    = regexp.MustCompile(`^[0-9]+\.[0-9]+`)
	major         = regexp.MustCompile(`^[0-9]+`)
)

// modelFamily maps a model id to its family name and rarity color.
func modelFamily(id string) (string, string) {
	for _, f := range []struct{ name, color string }{
		{"fable", colorDanger},
		{"mythos", colorDanger},
		{"opus", rarityLegendary},
		{"sonnet", rarityEpic},
		{"haiku", rarityRare},
	} {
		if strings.Contains(id, f.name) {
			return f.name, f.color
		}
	}
	return "unknown", colorFallback
}

// modelVersion strips "claude-" and the family name, then reads leading numeric
// segments as a version. e.g. claude-opus-4-7 → "4.7"; claude-3-5-sonnet-20241022 → "3.5".
func modelVersion(id, family string) string {
	s := strings.Trim(strings.Join(strings.Split(strings.ReplaceAll(id, "claude-", ""), family), ""), "-")
	segs := numericPrefix.FindString(s)
	if segs == "" {
		return ""
	}
	dotted := strings.ReplaceAll(segs, "-", ".")
	if v := majorMinor.FindString(dotted); v != "" {
		return v
	}
	return major.FindString(dotted)
}

// limitBar builds a rate-limit segment: {label}:{pct}% {bar} {H:MM until reset}.
func limitBar(label string, l limit) string {
	const barWidth = 10
	pct := int(math.Round(*l.UsedPercentage))
	filled := pct * barWidth / 100
	if filled < 0 {
		filled = 0
	} else if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("■", filled) + strings.Repeat("□", barWidth-filled)

	resetStr := ""
	if l.ResetsAt != nil {
		secsLeft := max(0, *l.ResetsAt-time.Now().Unix())
		resetStr = fmt.Sprintf(" %d:%02d", secsLeft/3600, secsLeft%3600/60)
	}
	return fmt.Sprintf("%s:%s%d%% %s%s%s", label, pctColor(pct), pct, bar, reset, resetStr)
}

func main() {
	// --- Read & parse input once ---
	var in input
	if raw, err := io.ReadAll(os.Stdin); err == nil {
		_ = json.Unmarshal(raw, &in) // malformed input → zero value, same as the JS fallback
	}

	var v view

	// --- Working Directory (home dir → ~), blue foreground ---
	cwd := in.Workspace.CurrentDir
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	dir := cwd
	if home := os.Getenv("HOME"); home != "" && strings.HasPrefix(cwd, home) {
		dir = "~" + cwd[len(home):]
	}
	v.Dir = blue + dir + reset

	// --- Git Branch & Diff Stats ---
	if branch := git(cwd, "rev-parse", "--abbrev-ref", "HEAD"); branch != "" {
		// Detached HEAD prints "HEAD"; fall back to short SHA to match prior behavior.
		label := branch
		if branch == "HEAD" {
			label = git(cwd, "rev-parse", "--short", "HEAD")
		}
		v.Branch = " on " + green + label + reset

		// Sum additions/deletions across all working-tree changes vs HEAD, skipping binaries ('-').
		// One `diff HEAD` scan replaces an unstaged+staged pair.
		var additions, deletions int
		for line := range strings.SplitSeq(git(cwd, "diff", "--numstat", "HEAD"), "\n") {
			cols := strings.Split(line, "\t")
			if len(cols) < 2 {
				continue
			}
			if n, err := strconv.Atoi(cols[0]); err == nil {
				additions += n
			}
			if n, err := strconv.Atoi(cols[1]); err == nil {
				deletions += n
			}
		}
		v.Diff = fmt.Sprintf("[%s+%d%s|%s-%d%s]", green, additions, reset, red, deletions, reset)
	}

	// --- Context Usage (color-coded) ---
	// Used tokens derived from used_percentage × context_window_size so the K value
	// stays consistent with the percentage Claude Code itself displays.
	usedPct, ctxWindow := in.ContextWindow.UsedPercentage, in.ContextWindow.ContextWindowSize

	sizeLabel := ""
	if ctxWindow != nil {
		if w := math.Round(*ctxWindow); w >= 1_000_000 {
			sizeLabel = fmt.Sprintf("%.0fM", math.Round(w/1_000_000))
		} else {
			sizeLabel = fmt.Sprintf("%.0fK", math.Round(w/1000))
		}
	}
	switch {
	case usedPct != nil && ctxWindow != nil:
		used := math.Round(*usedPct * *ctxWindow / 100 / 1000)
		v.Ctx = fmt.Sprintf("%s%.0fK%s/%s", pctColor(int(math.Round(*usedPct))), used, reset, sizeLabel)
	case usedPct != nil:
		pct := int(math.Round(*usedPct))
		v.Ctx = fmt.Sprintf("%s%d%%%s", pctColor(pct), pct, reset)
	case sizeLabel != "":
		v.Ctx = "?/" + sizeLabel
	}

	// --- Model & Effort (RPG-rarity color-coded by model family) ---
	if in.Model.ID != "" {
		id := strings.ToLower(in.Model.ID)
		family, color := modelFamily(id)

		shortName := family
		if version := modelVersion(id, family); version != "" {
			shortName = family + "-" + version
		}
		effort := in.Effort.Level
		if effort == "" {
			effort = "default"
		}
		v.Model = fmt.Sprintf("%s%s%s:%s%s%s", color, shortName, reset, color, effort, reset)
	}

	// --- Session cost (API-equivalent $, what these tokens would cost pay-as-you-go) ---
	for _, c := range []*float64{in.Cost.TotalCostUsd, in.Cost.TotalCostCC, in.Cost.TotalCost} {
		if c != nil {
			v.Cost = fmt.Sprintf("%s$%.2f%s", money, *c, reset)
			break
		}
	}

	// --- Rate limits ---
	if in.RateLimits.FiveHour.UsedPercentage != nil {
		v.Limits = append(v.Limits, limitBar("5h", in.RateLimits.FiveHour))
	}
	if in.RateLimits.SevenDay.UsedPercentage != nil {
		v.Limits = append(v.Limits, limitBar("7d", in.RateLimits.SevenDay))
	}

	template.Must(template.New("statusline").Parse(layout)).Execute(os.Stdout, v)
}
