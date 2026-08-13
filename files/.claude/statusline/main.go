// Claude Code status line. Build the binary with `go build` in this directory.
// The program reads the session JSON from stdin and runs git for the branch and
// the diff counts. It prints two lines and omits the empty segments:
//
//	{directory} on {branch} [+N|-N] | {cost}
//	{model}[{context size}]:{effort} | {tokens} {bar} | {5h limit} {bar} | {7d limit} {bar}
//
// Each limit segment shows the time until the limit resets.
package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// --- ANSI colors ---
const (
	reset  = "\x1b[0m"
	green  = "\x1b[32m"
	yello  = "\x1b[33m"
	orange = "\x1b[38;5;208m"
	red    = "\x1b[31m"
	blue   = "\x1b[34m"
	money  = "\x1b[38;5;78m" // soft green for the cost in dollars

	// RPG rarity colors (256-color) for the model families.
	// Rare is Haiku, Epic is Sonnet, Legendary is Opus.
	// Danger (crimson) marks Fable and Mythos, the high-cost models.
	rarityRare      = "\x1b[38;5;69m"
	rarityEpic      = "\x1b[38;5;135m"
	rarityLegendary = "\x1b[38;5;220m"
	colorDanger     = "\x1b[38;2;220;20;60m" // crimson #DC143C, a cost warning
	colorFallback   = "\x1b[38;5;246m"       // unknown model
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
	// The field name changed between Claude Code versions. Read all three spellings.
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

// view holds the segments of both lines. Each segment carries its own color.
type view struct {
	Dir     string
	Branch  string
	Diff    string
	Model   string
	Ctx     string
	Cost    string
	Limit5h string
	Limit7d string
}

// layout places the segments. gitStatus sets Branch and Diff together, so one
// guard covers both. The join function drops the empty segments of line 2.
const layout = `{{.Dir}}{{with .Branch}} on {{.}} {{$.Diff}}{{end}}{{with .Cost}} | {{.}}{{end}}
{{join .Model .Ctx .Limit5h .Limit7d}}
`

var statusLine = template.Must(template.New("statusline").
	Funcs(template.FuncMap{"join": join}).Parse(layout))

// join concatenates the segments that are not empty, separated by " | ".
func join(segments ...string) string {
	return strings.Join(slices.DeleteFunc(segments, func(s string) bool { return s == "" }), " | ")
}

// git runs one git command. It returns an empty string after any failure, for
// example a directory outside a repository or a missing git binary.
func git(cwd string, args ...string) string {
	out, err := exec.Command("git", append([]string{"-C", cwd, "--no-optional-locks"}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// gitStatus returns the colored branch name and the line counts of the working
// tree. Both are empty when cwd is outside a repository.
func gitStatus(cwd string) (branch, diff string) {
	// A git process costs about 12 ms and the two commands are independent.
	// The scan starts first and runs while the branch name arrives.
	// One `diff HEAD` scan covers both the staged and the unstaged changes.
	numstat := make(chan string, 1)
	go func() { numstat <- git(cwd, "diff", "--numstat", "HEAD") }()

	name := git(cwd, "rev-parse", "--abbrev-ref", "HEAD")
	if name == "" {
		return "", ""
	}
	if name == "HEAD" { // a detached HEAD has no name, so show the short hash
		name = git(cwd, "rev-parse", "--short", "HEAD")
	}

	// A binary file reports "-". Atoi then fails and the count stays the same.
	var additions, deletions int
	for line := range strings.SplitSeq(<-numstat, "\n") {
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
	return green + name + reset,
		fmt.Sprintf("[%s+%d%s|%s-%d%s]", green, additions, reset, red, deletions, reset)
}

// pctColor colors a percentage: 40 or less green, 50 or less yellow,
// 80 or less orange, above 80 red.
func pctColor(pct int) string {
	switch {
	case pct <= 40:
		return green
	case pct <= 50:
		return yello
	case pct <= 80:
		return orange
	default:
		return red
	}
}

// gauge renders "{value} {bar}". The percentage sets the color and the fill.
// A nil percentage gives an empty segment, and an empty value gives "{pct}%".
func gauge(usedPct *float64, value string) string {
	const barWidth = 10
	if usedPct == nil {
		return ""
	}
	pct := int(math.Round(*usedPct))
	if value == "" {
		value = fmt.Sprintf("%d%%", pct)
	}
	filled := min(max(pct*barWidth/100, 0), barWidth)
	bar := strings.Repeat("■", filled) + strings.Repeat("□", barWidth-filled)
	return pctColor(pct) + value + " " + bar + reset
}

var numericPrefix = regexp.MustCompile(`^([0-9]+)(?:-([0-9]+))?`)

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

// modelVersion removes "claude-" and the family name, then reads the first two
// numeric segments as a version. claude-opus-4-7 gives "4.7", and
// claude-3-5-sonnet-20241022 gives "3.5".
func modelVersion(id, family string) string {
	s := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(id, "claude-", ""), family, ""), "-")
	m := numericPrefix.FindStringSubmatch(s)
	switch {
	case m == nil:
		return ""
	case m[2] != "":
		return m[1] + "." + m[2]
	}
	return m[1]
}

// modelSegment renders "{family}-{version}[{context size}]:{effort}" in the
// family color, for example "opus-5[1M]:xhigh".
func modelSegment(id, effort string, window *float64) string {
	if id == "" {
		return ""
	}
	id = strings.ToLower(id)
	family, color := modelFamily(id)
	name := family
	if version := modelVersion(id, family); version != "" {
		name = family + "-" + version
	}
	if window != nil {
		name += "[" + tokens(*window) + "]"
	}
	if effort == "" {
		effort = "default"
	}
	return color + name + reset + ":" + color + effort + reset
}

// tokens formats a token count, for example "1M" or "200K".
func tokens(n float64) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.0fM", n/1_000_000)
	}
	return fmt.Sprintf("%.0fK", n/1000)
}

// contextGauge shows the tokens in use, for example "70K ■□□□□□□□□□". The token
// count comes from the percentage, so it agrees with the number Claude Code shows.
func contextGauge(usedPct, window *float64) string {
	value := ""
	if usedPct != nil && window != nil {
		value = tokens(*usedPct * *window / 100)
	}
	return gauge(usedPct, value)
}

// limitGauge shows the time until the limit resets, next to the usage bar.
// The value is "1:39" below one day and "5d 20h" above it.
func limitGauge(l limit) string {
	value := ""
	if l.ResetsAt != nil {
		secsLeft := max(0, *l.ResetsAt-time.Now().Unix())
		if hours := secsLeft / 3600; hours >= 24 {
			value = fmt.Sprintf("%dd %dh", hours/24, hours%24)
		} else {
			value = fmt.Sprintf("%d:%02d", hours, secsLeft%3600/60)
		}
	}
	return gauge(l.UsedPercentage, value)
}

// costSegment formats the first cost that the session reports. The value is the
// API price of these tokens.
func costSegment(costs ...*float64) string {
	c := cmp.Or(costs...)
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%s$%.2f%s", money, *c, reset)
}

// dirSegment shortens the home directory to "~" and colors the path.
func dirSegment(dir string) string {
	if home, err := os.UserHomeDir(); err == nil {
		if rest, found := strings.CutPrefix(dir, home); found {
			dir = "~" + rest
		}
	}
	return blue + dir + reset
}

func main() {
	// Stop early when stdin is a terminal. A direct `go run` carries no session JSON.
	if fi, err := os.Stdin.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintln(os.Stderr, "statusline: no input on stdin; Claude Code must run this command")
		os.Exit(1)
	}

	// Empty input ends the decode with io.EOF, which fails here too.
	var in input
	if json.NewDecoder(os.Stdin).Decode(&in) != nil || (in.Workspace.CurrentDir == "" && in.Model.ID == "") {
		fmt.Fprintln(os.Stderr, "statusline: no valid Claude Code session data on stdin")
		os.Exit(1)
	}

	cwd := in.Workspace.CurrentDir
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	branch, diff := gitStatus(cwd)
	ctx := in.ContextWindow

	statusLine.Execute(os.Stdout, view{
		Dir:     dirSegment(cwd),
		Branch:  branch,
		Diff:    diff,
		Model:   modelSegment(in.Model.ID, in.Effort.Level, ctx.ContextWindowSize),
		Ctx:     contextGauge(ctx.UsedPercentage, ctx.ContextWindowSize),
		Cost:    costSegment(in.Cost.TotalCostUsd, in.Cost.TotalCostCC, in.Cost.TotalCost),
		Limit5h: limitGauge(in.RateLimits.FiveHour),
		Limit7d: limitGauge(in.RateLimits.SevenDay),
	})
}
