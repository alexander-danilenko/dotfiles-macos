#!/usr/bin/env node
// Claude Code status line — Node.js rewrite for speed.
// Replaces the previous bash script that spawned ~11 jq calls plus many
// sed/grep/tr/awk subshells per render. This parses the JSON once and does all
// string work in-process, shelling out only for git.
// Template:
//   Line 1: {{pwd blue}} on {{git_branch green}} [+N|-N]
//   Line 2: {model}:{effort} | ctx:{used}/{total} [| 5h:{pct}% bar][| 7d:{pct}% bar]
"use strict";

const fs = require("fs");
const { execFileSync } = require("child_process");

// --- ANSI helpers ---
const reset = "\x1b[0m";
const green = "\x1b[32m";
const yellow = "\x1b[33m";
const red = "\x1b[31m";
const blue = "\x1b[34m";
const money = "\x1b[38;5;78m"; // soft green for the API-equivalent $ cost

// --- RPG rarity colors (256-color) for model families ---
// Rare → Haiku, Epic → Sonnet, Legendary → Opus, Mythic → Mythos (reserved).
const rarityRare = "\x1b[38;5;69m"; // Haiku
const rarityEpic = "\x1b[38;5;135m"; // Sonnet
const rarityLegendary = "\x1b[38;5;220m"; // Opus
const rarityMythic = "\x1b[38;5;159m"; // Mythos (reserved)
const colorFallback = "\x1b[38;5;246m"; // unknown/fallback

// --- Read & parse input once ---
let input = {};
try {
  input = JSON.parse(fs.readFileSync(0, "utf8") || "{}");
} catch {
  input = {};
}

const get = (obj, path) =>
  path.split(".").reduce((acc, k) => (acc == null ? undefined : acc[k]), obj);

// Run git quietly; return '' on any failure (not a repo, git missing, etc.).
const git = (cwd, args) => {
  try {
    return execFileSync("git", ["-C", cwd, "--no-optional-locks", ...args], {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
    }).trim();
  } catch {
    return "";
  }
};

// Color for a percentage: <=60 green, <=80 yellow, else red.
const pctColor = (pct) => (pct <= 60 ? green : pct <= 80 ? yellow : red);

// --- Working Directory (home dir → ~), blue foreground ---
const cwd = get(input, "workspace.current_dir") || process.cwd();
const home = process.env.HOME || "";
const shortenedRaw =
  home && cwd.startsWith(home) ? "~" + cwd.slice(home.length) : cwd;
const shortenedDir = `${blue}${shortenedRaw}${reset}`;

// --- Git Branch & Diff Stats ---
let gitBranchInfo = "";
let gitDiffInfo = "";
const branch = git(cwd, ["rev-parse", "--abbrev-ref", "HEAD"]);
if (branch) {
  // Detached HEAD prints "HEAD"; fall back to short SHA to match prior behavior.
  const label =
    branch === "HEAD" ? git(cwd, ["rev-parse", "--short", "HEAD"]) : branch;
  gitBranchInfo = ` on ${green}${label}${reset}`;

  // Sum additions/deletions across all working-tree changes vs HEAD, skipping binaries ('-').
  // One `diff HEAD` scan replaces the old unstaged+staged pair (verified equal for clean and
  // dirty trees) — this is the single change that makes the status line measurably faster.
  const numstat = git(cwd, ["diff", "--numstat", "HEAD"]);
  let additions = 0;
  let deletions = 0;
  for (const line of numstat.split("\n")) {
    const [add, del] = line.split("\t");
    if (/^\d+$/.test(add)) additions += Number(add);
    if (/^\d+$/.test(del)) deletions += Number(del);
  }
  gitDiffInfo = `[${green}+${additions}${reset}|${red}-${deletions}${reset}]`;
}

// --- Context Usage (color-coded) ---
// Used tokens derived from used_percentage × context_window_size so the K value
// stays consistent with the percentage Claude Code itself displays.
const usedPct = get(input, "context_window.used_percentage");
const ctxWindow = get(input, "context_window.context_window_size");

let ctxSizeLabel = "";
if (ctxWindow != null) {
  const w = Math.round(ctxWindow);
  ctxSizeLabel =
    w >= 1000000 ? `${Math.round(w / 1000000)}M` : `${Math.round(w / 1000)}K`;
}

let ctxTokensStr = "";
if (usedPct != null) {
  const pctInt = Math.round(usedPct);
  const color = pctColor(pctInt);
  if (ctxWindow != null) {
    const usedLabel = `${Math.round((usedPct * ctxWindow) / 100 / 1000)}K`;
    ctxTokensStr = `${color}${usedLabel}${reset}/${ctxSizeLabel}`;
  } else {
    ctxTokensStr = `${color}${pctInt}%${reset}`;
  }
}

// --- Model & Effort (RPG-rarity color-coded by model family) ---
let modelInfo = "";
const modelId = get(input, "model.id");
let effort = get(input, "effort.level");

if (modelId) {
  const modelLower = String(modelId).toLowerCase();

  let family;
  let modelColor;
  if (modelLower.includes("mythos")) {
    family = "mythos";
    modelColor = rarityMythic;
  } else if (modelLower.includes("opus")) {
    family = "opus";
    modelColor = rarityLegendary;
  } else if (modelLower.includes("sonnet")) {
    family = "sonnet";
    modelColor = rarityEpic;
  } else if (modelLower.includes("haiku")) {
    family = "haiku";
    modelColor = rarityRare;
  } else {
    family = "unknown";
    modelColor = colorFallback;
  }

  // Strip "claude-" and the family name, then read leading numeric segments as a version.
  // e.g. claude-opus-4-7 → "4.7"; claude-3-5-sonnet-20241022 → "3.5".
  const stripped = modelLower
    .replace(/claude-/g, "")
    .split(family)
    .join("")
    .replace(/^-/, "")
    .replace(/-$/, "");
  const segs = stripped.match(/^[0-9]+(-[0-9]+)*/);
  let version = "";
  if (segs) {
    const dotted = segs[0].replace(/-/g, ".");
    const dotVer = dotted.match(/^[0-9]+\.[0-9]+/);
    version = dotVer ? dotVer[0] : dotted.match(/^[0-9]+/)[0];
  }

  const shortName = version ? `${family}-${version}` : family;
  if (!effort) effort = "default";

  modelInfo = `${modelColor}${shortName}${reset}:${modelColor}${effort}${reset}`;
}

// --- Second line: {model}:{effort} | ctx:{used}/{total} | 5h:{pct}% {bar} | 7d:{pct}% {bar} ---
// Build a rate-limit segment: {label}:{pct}% {bar} {Xh Ym}.
const limitBar = (label, pctFloat, resetsAt) => {
  const barWidth = 10;
  const pctInt = Math.round(pctFloat);
  const color = pctColor(pctInt);
  const filled = Math.floor((pctInt * barWidth) / 100);
  const bar = "■".repeat(filled) + "□".repeat(barWidth - filled);

  let resetStr = "";
  if (resetsAt != null) {
    const secsLeft = Math.max(0, resetsAt - Math.floor(Date.now() / 1000));
    const hrs = secsLeft / 3600;
    resetStr = ` ${hrs.toFixed(2)}h`;
  }

  return `${label}:${color}${pctInt}% ${bar}${reset}${resetStr}`;
};

// --- Session cost (API-equivalent $, what these tokens would cost pay-as-you-go) ---
// Claude Code precomputes this; field name has varied across versions, so probe
// the known spellings and fall back to nothing if absent.
const totalCostUsd =
  get(input, "cost.total_cost_usd") ??
  get(input, "cost.totalCost") ??
  get(input, "cost.total_cost") ??
  null;

let costStr = "";
if (totalCostUsd != null && !Number.isNaN(Number(totalCostUsd))) {
  const c = Number(totalCostUsd);
  costStr = `${money}$${c.toFixed(2)}${reset}`;
}

let secondLine = "";
if (modelInfo) {
  secondLine = modelInfo;

  if (ctxTokensStr) {
    secondLine += ` | ctx:${ctxTokensStr}`;
  } else if (ctxSizeLabel) {
    secondLine += ` | ctx:?/${ctxSizeLabel}`;
  }

  if (costStr) {
    secondLine += ` | ${costStr}`;
  }

  const fivePct = get(input, "rate_limits.five_hour.used_percentage");
  if (fivePct != null) {
    secondLine += ` | ${limitBar("5h", fivePct, get(input, "rate_limits.five_hour.resets_at"))}`;
  }
  const weekPct = get(input, "rate_limits.seven_day.used_percentage");
  if (weekPct != null) {
    secondLine += ` | ${limitBar("7d", weekPct, get(input, "rate_limits.seven_day.resets_at"))}`;
  }
}

// --- Assemble output ---
const line1 = `${shortenedDir}${gitBranchInfo} ${gitDiffInfo}`;
process.stdout.write(secondLine ? `${line1}\n${secondLine}` : line1);
