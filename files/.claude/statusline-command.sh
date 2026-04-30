#!/usr/bin/env bash
# Claude Code status line script
# ~/.claude/statusline-command.sh
# Template: {{pwd blue}} on {{git_branch green}} [+N|-N][model:effort usedK/size]

input=$(cat)

# --- ANSI helpers ---
reset='\033[0m'
bold='\033[1m'
green='\033[32m'
orange='\033[33m'
yellow='\033[33m'
red='\033[31m'
blue='\033[34m'
bright_blue='\033[1;94m'
purple='\033[35m'

# --- RPG rarity colors (256-color) for model families ---
# Rare     → Haiku   (steel blue)
# Epic     → Sonnet  (bright purple)
# Legendary→ Opus    (bright gold)
# Mythic   → Mythos  (icy cyan-white — placeholder for future model)
rarity_rare='\033[38;5;69m'       # Haiku    — Rare
rarity_epic='\033[38;5;135m'      # Sonnet   — Epic
rarity_legendary='\033[38;5;220m' # Opus     — Legendary
rarity_mythic='\033[38;5;159m'    # Mythos   — Mythic (reserved, not yet in production)
color_fallback='\033[38;5;246m'   # unknown/fallback model — neutral grey

# --- Working Directory (home dir replaced with ~/), blue foreground ---
cwd=$(echo "$input" | jq -r '.workspace.current_dir')
home_dir="$HOME"
shortened_raw=$(echo "$cwd" | sed "s|^$home_dir|~|")
shortened_dir=$(printf "${blue}%s${reset}" "$shortened_raw")

# --- Git Branch & Diff Stats ---
git_branch_info=""
git_diff_info=""
if git -C "$cwd" rev-parse --git-dir > /dev/null 2>&1; then
  branch=$(git -C "$cwd" --no-optional-locks symbolic-ref --short HEAD 2>/dev/null \
    || git -C "$cwd" --no-optional-locks rev-parse --short HEAD 2>/dev/null)
  # Green foreground
  git_branch_info=" on $(printf "${green}%s${reset}" "$branch")"

  # Diff stats: sum all additions and deletions across changed files (unstaged + staged)
  diff_raw=$(git -C "$cwd" --no-optional-locks diff --numstat 2>/dev/null)
  staged_raw=$(git -C "$cwd" --no-optional-locks diff --numstat --cached 2>/dev/null)
  combined_raw=$(printf '%s\n%s\n' "$diff_raw" "$staged_raw")

  additions=0
  deletions=0
  while IFS=$'\t' read -r add del _rest; do
    # Skip binary files (shown as '-')
    [[ "$add" =~ ^[0-9]+$ ]] && additions=$((additions + add))
    [[ "$del" =~ ^[0-9]+$ ]] && deletions=$((deletions + del))
  done <<< "$combined_raw"

  # Always show diff block in a git repo (show zeros if no changes)
  git_diff_info=$(printf "[${green}+%d${reset}|${red}-%d${reset}]" "$additions" "$deletions")
fi

# --- Context Usage (color-coded) ---
# <=60%: green  61-80%: yellow  >80%: red
context_info=""
used_pct=$(echo "$input" | jq -r '.context_window.used_percentage // empty')
ctx_window=$(echo "$input" | jq -r '.context_window.context_window_size // empty')
ctx_input_tokens=$(echo "$input" | jq -r '.context_window.current_usage.input_tokens // empty')

# Format context window size as K or M
ctx_size_label=""
if [ -n "$ctx_window" ] && [ "$ctx_window" != "null" ]; then
  ctx_window_int=$(printf "%.0f" "$ctx_window")
  if [ "$ctx_window_int" -ge 1000000 ]; then
    ctx_size_label=$(awk "BEGIN { printf \"%.0fM\", $ctx_window_int/1000000 }")
  else
    ctx_size_label=$(awk "BEGIN { printf \"%.0fK\", $ctx_window_int/1000 }")
  fi
fi

# Format used tokens as rounded K value (e.g. 123456 → 123K)
ctx_used_label=""
if [ -n "$ctx_input_tokens" ] && [ "$ctx_input_tokens" != "null" ]; then
  ctx_used_label=$(awk "BEGIN { printf \"%.0fK\", $ctx_input_tokens/1000 }")
fi

if [ -n "$used_pct" ]; then
  pct_int=$(printf "%.0f" "$used_pct")
  if [ "$pct_int" -le 60 ]; then
    ctx_color="$green"    # green — plenty of room
  elif [ "$pct_int" -le 80 ]; then
    ctx_color="$yellow"   # yellow — getting full
  else
    ctx_color="$red"      # red — nearly exhausted
  fi

  # Build the used/size portion — use printf so \033 escapes in color vars are interpreted
  if [ -n "$ctx_used_label" ] && [ -n "$ctx_size_label" ]; then
    ctx_tokens_str=$(printf "${ctx_color}%s${reset}/%s" "$ctx_used_label" "$ctx_size_label")
  elif [ -n "$ctx_used_label" ]; then
    ctx_tokens_str=$(printf "${ctx_color}%s${reset}" "$ctx_used_label")
  elif [ -n "$ctx_size_label" ]; then
    ctx_tokens_str=$(printf "${ctx_color}%d%%${reset}/%s" "$pct_int" "$ctx_size_label")
  else
    ctx_tokens_str=$(printf "${ctx_color}%d%%${reset}" "$pct_int")
  fi
fi

# --- Model & Effort (RPG-rarity color-coded by model family) ---
# haiku=Rare(blue)  sonnet=Epic(purple)  opus=Legendary(gold)
# mythos=Mythic(cyan-white) — placeholder for future release
# other/unknown → plain purple fallback
model_info=""
model_id=$(echo "$input" | jq -r '.model.id // empty')
effort=$(echo "$input" | jq -r '.effort.level // empty')

if [ -n "$model_id" ]; then
  model_lower=$(echo "$model_id" | tr '[:upper:]' '[:lower:]')

  # Detect family and assign RPG rarity color
  if echo "$model_lower" | grep -q "mythos"; then
    # --- MYTHIC TIER (future model: claude-mythos-*) ---
    model_family="mythos"
    model_color="$rarity_mythic"
  elif echo "$model_lower" | grep -q "opus"; then
    model_family="opus"
    model_color="$rarity_legendary"
  elif echo "$model_lower" | grep -q "sonnet"; then
    model_family="sonnet"
    model_color="$rarity_epic"
  elif echo "$model_lower" | grep -q "haiku"; then
    model_family="haiku"
    model_color="$rarity_rare"
  else
    model_family="unknown"
    model_color="$color_fallback"
  fi

  # Extract version numbers from model ID
  # e.g. claude-sonnet-4-6 → 4.6, claude-opus-4-5 → 4.5, claude-3-5-sonnet-20241022 → 3.5
  # Strategy: strip the family name and "claude-" prefix, then grab leading digit groups
  stripped=$(echo "$model_lower" | sed "s/claude-//g" | sed "s/$model_family//g" | sed 's/^-//' | sed 's/-$//')
  # Extract leading numeric segments (e.g. "4-6-..." → "4.6", "3-5-..." → "3.5")
  version=$(echo "$stripped" | grep -oE '^[0-9]+(-[0-9]+)*' | sed 's/-/./g' | grep -oE '^[0-9]+\.[0-9]+')
  if [ -z "$version" ]; then
    # Fallback: just the first number if no dot-version found
    version=$(echo "$stripped" | grep -oE '^[0-9]+')
  fi

  if [ -n "$version" ]; then
    short_name="${model_family}-${version}"
  else
    short_name="$model_family"
  fi

  # Effort: ALWAYS show — even when "default" or "auto" or empty
  if [ -z "$effort" ]; then
    effort="default"
  fi

  model_info=$(printf "${model_color}%s${reset}:${model_color}%s${reset}" "$short_name" "$effort")
fi

# --- Combine model+effort with context into a single bracket ---
# Format: [model:effort usedK/size]
if [ -n "$model_info" ] && [ -n "$used_pct" ]; then
  context_info="[${model_info} ${ctx_tokens_str}]"
elif [ -n "$model_info" ]; then
  context_info="[${model_info}]"
elif [ -n "$used_pct" ]; then
  context_info="[${ctx_tokens_str}]"
fi

# --- Usage / Rate Limits (5-hour and 7-day) ---
# Compact inline bracket appended to the main line.
# Only rendered when at least one limit is present in the JSON.
# Format: [5h:<pct>%|7d:<pct>%]  — each token colored by its own threshold.
# Color thresholds: <=60% green, <=80% yellow, >80% red.
limits_info=""
five_pct=$(echo "$input" | jq -r '.rate_limits.five_hour.used_percentage // empty')
week_pct=$(echo "$input" | jq -r '.rate_limits.seven_day.used_percentage // empty')

limit_token() {
  local label="$1"
  local pct_float="$2"
  local pct_int
  pct_int=$(printf "%.0f" "$pct_float")
  local tok_color
  if [ "$pct_int" -le 60 ]; then
    tok_color="$green"
  elif [ "$pct_int" -le 80 ]; then
    tok_color="$yellow"
  else
    tok_color="$red"
  fi
  printf "${tok_color}%s:%d%%${reset}" "$label" "$pct_int"
}

if [ -n "$five_pct" ] || [ -n "$week_pct" ]; then
  limit_parts=""
  if [ -n "$five_pct" ]; then
    limit_parts="$(limit_token "5h" "$five_pct")"
  fi
  if [ -n "$week_pct" ]; then
    tok_7d="$(limit_token "7d" "$week_pct")"
    if [ -n "$limit_parts" ]; then
      limit_parts="${limit_parts}|${tok_7d}"
    else
      limit_parts="$tok_7d"
    fi
  fi
  limits_info="[${limit_parts}]"
fi

# --- Assemble output ---
# Single line: {{pwd blue}} on {{git_branch green}} [+N|-N][model:effort usedK/size][5h:N%|7d:N%]
printf "%s%s %s%s%s" \
  "$shortened_dir" \
  "$git_branch_info" \
  "$git_diff_info" \
  "$context_info" \
  "$limits_info"
