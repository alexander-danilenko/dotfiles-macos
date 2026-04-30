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
# Used tokens are derived from used_percentage × context_window_size so the K value
# stays consistent with the percentage Claude Code itself displays. Reading
# .current_usage.input_tokens directly under-counts because prompt caching moves
# most input volume into cache_read_input_tokens / cache_creation_input_tokens.
context_info=""
used_pct=$(echo "$input" | jq -r '.context_window.used_percentage // empty')
ctx_window=$(echo "$input" | jq -r '.context_window.context_window_size // empty')

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

# Used tokens: used_pct × ctx_window / 100, rounded to K (e.g. 8% × 200000 → 16K).
ctx_used_label=""
if [ -n "$used_pct" ] && [ -n "$ctx_window" ] && [ "$ctx_window" != "null" ]; then
  ctx_used_label=$(awk "BEGIN { printf \"%.0fK\", ($used_pct * $ctx_window / 100) / 1000 }")
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

# --- Second line: {model}:{effort} | ctx:{used}/{total} | 5h:{pct}% {bar} | 7d:{pct}% {bar} ---
# ctx segment always shown when model info is present.
# Rate-limit segments only shown when data is available.
# Color thresholds: <=60% green, <=80% yellow, >80% red.
# Bar width: 20 characters.
second_line=""
five_pct=$(echo "$input" | jq -r '.rate_limits.five_hour.used_percentage // empty')
week_pct=$(echo "$input" | jq -r '.rate_limits.seven_day.used_percentage // empty')

# Build a rate-limit segment: {label}:{pct}% {bar} {Xh Ym}
# Args: $1=label (e.g. "5h"), $2=pct_float, $3=resets_at (Unix epoch, optional)
limit_bar() {
  local label="$1"
  local pct_float="$2"
  local resets_at="${3:-}"
  local bar_width=10
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
  local filled=$(( pct_int * bar_width / 100 ))
  local empty=$(( bar_width - filled ))
  local bar=""
  local i
  for (( i=0; i<filled; i++ )); do bar="${bar}■"; done
  for (( i=0; i<empty;  i++ )); do bar="${bar}□"; done

  # Reset time: compute Xh Ym remaining from resets_at epoch
  local reset_str=""
  if [ -n "$resets_at" ] && [ "$resets_at" != "null" ]; then
    local now
    now=$(date +%s)
    local secs_left=$(( resets_at - now ))
    if [ "$secs_left" -lt 0 ]; then
      secs_left=0
    fi
    local hrs=$(( secs_left / 3600 ))
    local mins=$(( (secs_left % 3600) / 60 ))
    reset_str=$(printf " %dh %dm" "$hrs" "$mins")
  fi

  printf "%s:${tok_color}%d%% %s${reset}%s" "$label" "$pct_int" "$bar" "$reset_str"
}

if [ -n "$model_info" ]; then
  # Start with model:effort
  second_line="$model_info"

  # ctx segment: always append when model is known
  if [ -n "$ctx_tokens_str" ]; then
    second_line="${second_line} | ctx:${ctx_tokens_str}"
  elif [ -n "$ctx_size_label" ]; then
    second_line="${second_line} | ctx:?/${ctx_size_label}"
  fi

  # Rate-limit segments: only when data is present
  if [ -n "$five_pct" ]; then
    five_resets=$(echo "$input" | jq -r '.rate_limits.five_hour.resets_at // empty')
    second_line="${second_line} | $(limit_bar "5h" "$five_pct" "$five_resets")"
  fi
  if [ -n "$week_pct" ]; then
    week_resets=$(echo "$input" | jq -r '.rate_limits.seven_day.resets_at // empty')
    second_line="${second_line} | $(limit_bar "7d" "$week_pct" "$week_resets")"
  fi
fi

# --- Assemble output ---
# Line 1: {{pwd blue}} on {{git_branch green}} [+N|-N]
# Line 2 (when model is known): {model}:{effort} | ctx:{used}/{total} [| 5h:{pct}% bar][| 7d:{pct}% bar]
if [ -n "$second_line" ]; then
  printf "%s%s %s\n%s" \
    "$shortened_dir" \
    "$git_branch_info" \
    "$git_diff_info" \
    "$second_line"
else
  printf "%s%s %s" \
    "$shortened_dir" \
    "$git_branch_info" \
    "$git_diff_info"
fi
