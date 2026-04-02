# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## What This Is

This is the configuration directory for [AeroSpace](https://github.com/nikitabobko/AeroSpace), an i3-like tiling window manager for macOS. The entire config lives in a single file: `aerospace.toml` (TOML format, `config-version = 2`).

Always consult AeroSpace documentation (via Context7 MCP with library ID `/nikitabobko/aerospace`) before answering questions or making changes, as AeroSpace is actively developed and options may change between versions.

## Config Architecture

This config follows a **fixed app-to-workspace-to-monitor** paradigm:

- **7 persistent workspaces** with designated purposes:
  - Main monitor (1-4): Chat (Slack/Telegram), Mail/Calendar, JetBrains IDEs, Zed
  - Secondary monitor (5-7): Terminal (Ghostty), Chrome, Vivaldi
- **`on-window-detected` callbacks** auto-assign apps to their workspace by `app-id` or `app-name-regex-substring`
- **`workspace-to-monitor-force-assignment`** pins workspaces to monitors using `'main'`/`'secondary'` keywords
- **`move-node-to-workspace` bindings are intentionally omitted** to prevent apps from leaving their assigned workspaces — do not add them

## Key Design Decisions

- All windows default to tiling layout via a catch-all `on-window-detected` with `check-further-callbacks = true`
- JetBrains IDEs are matched by name regex (not app-id) to cover all IDE variants
- Mail and Calendar share workspace 2 with explicit positioning (`move left`/`move right`) and Calendar resized narrower
- Two modes: **main** (normal operation, vim-style navigation with alt+hjkl) and **service** (entered via `alt-shift-;`, exited with `esc` which also reloads config)

## Useful Commands

```bash
# Validate config (AeroSpace reloads on save, or use service mode: alt-shift-; then esc)
aerospace reload-config

# Debug window detection — find app-id for a running app
aerospace list-apps

# List current workspace tree
aerospace list-windows --all
```

## Editing Guidelines

- **After every config change, run `aerospace reload-config`** to apply it immediately
- When adding a new app assignment, use `app-id` (from `aerospace list-apps`) unless multiple app variants need matching, then use `app-name-regex-substring`
- PWA app-ids are opaque hashes (e.g. `com.vivaldi.Vivaldi.app.gdfaincndogidkdcdkhapmbffkckdkhn`) — always add a comment above the `[[on-window-detected]]` block with the actual app name
- Place new `[[on-window-detected]]` blocks in the correct workspace section, following the existing comment structure
- Keep the workspace numbering and purpose comments in sync across: the top comment block, `persistent-workspaces`, monitor assignment, `on-window-detected` sections, and keybinding comments
- Maintain comprehensive comments in `aerospace.toml` — document the reasoning behind each section, workspace purpose, and non-obvious choices so the config remains self-explanatory
