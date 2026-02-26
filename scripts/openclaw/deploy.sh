#!/usr/bin/env bash
set -euo pipefail

mode="${1:-dry-run}"
if [[ "$mode" != "dry-run" ]]; then
  echo "only dry-run mode is supported in this script" >&2
  exit 2
fi

required_files=(
  "config/openclaw/SKILL.md"
  "config/openclaw/telegram.example.jsonc"
  "config/openclaw/mcp-server.example.jsonc"
  "scripts/openclaw/openppl-automation.sh"
  "scripts/openclaw/smoke.sh"
)

for f in "${required_files[@]}"; do
  [[ -f "$f" ]] || { echo "missing required file: $f" >&2; exit 1; }
done

grep -q "status" config/openclaw/SKILL.md
grep -q "send reminder" config/openclaw/SKILL.md
grep -q "deny" config/openclaw/SKILL.md
grep -q "dmPolicy" config/openclaw/telegram.example.jsonc
grep -q "pairing" config/openclaw/telegram.example.jsonc
grep -q "requireMention" config/openclaw/telegram.example.jsonc
grep -q "openppl-automation.sh" config/openclaw/mcp-server.example.jsonc

echo "dry-run checks passed"
