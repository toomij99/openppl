# OpenClaw + openppl Integration Guide

This guide shows how to connect your existing OpenClaw Telegram bot to `openppl` using the built-in bounded skill and wrapper scripts.

## What You Get

- `status` intent -> returns `openppl automation status` JSON
- `send reminder` intent -> calls `openppl automation action --name remind ...`
- deny-by-default behavior for unknown actions
- pairing-first Telegram defaults

## Prerequisites

- OpenClaw is already running with your Telegram bot
- `openppl` is installed and runnable from shell (`openppl` command on `PATH`)
- You are in the `openppl` repo root when running the scripts below

If `openppl` is not on your `PATH`, set it explicitly:

```bash
export OPENPPL_BIN="$(pwd)/openppl"
```

## Included Files (in this repo)

- `config/openclaw/SKILL.md`
- `config/openclaw/telegram.example.jsonc`
- `config/openclaw/mcp-server.example.jsonc`
- `scripts/openclaw/openppl-automation.sh`
- `scripts/openclaw/deploy.sh`
- `scripts/openclaw/smoke.sh`

## Quick Install Flow (Recommended)

Run this sequence from the `openppl` repo root:

```bash
# 1) Validate assets and config templates
bash scripts/openclaw/deploy.sh dry-run

# 2) Verify openppl command contract
openppl automation status
openppl automation action --name remind --request-id setup-001 --actor-scope telegram:setup

# 3) Verify wrapper contract
bash scripts/openclaw/openppl-automation.sh status
bash scripts/openclaw/openppl-automation.sh remind --request-id setup-002 --actor-scope telegram:setup

# 4) Run smoke tests
bash tests/openclaw/smoke_test.sh
```

If all commands pass, proceed to OpenClaw runtime config and skill registration.

## One-Shot Prompt for OpenClaw

Copy/paste this prompt into OpenClaw to set up openppl skill wiring end-to-end:

```text
You are configuring OpenClaw integration for openppl in this repository.

Goals:
1) Install/register a bounded openppl skill for Telegram.
2) Allow only two actions: status and remind.
3) Deny unknown actions by default.
4) Keep pairing-first policy and mention requirement in groups.
5) Validate with wrapper + smoke tests.

Execution constraints:
- Work from repository root.
- Never expose or write real secrets to git-tracked files.
- Do not enable arbitrary shell passthrough.
- Do not add actions beyond status/remind.

Source files to use:
- config/openclaw/SKILL.md
- config/openclaw/telegram.example.jsonc
- config/openclaw/mcp-server.example.jsonc
- scripts/openclaw/openppl-automation.sh
- scripts/openclaw/deploy.sh
- tests/openclaw/smoke_test.sh

Required steps:
1) Run: bash scripts/openclaw/deploy.sh dry-run
2) Verify commands:
   - openppl automation status
   - openppl automation action --name remind --request-id setup-001 --actor-scope telegram:setup
3) Verify wrapper:
   - bash scripts/openclaw/openppl-automation.sh status
   - bash scripts/openclaw/openppl-automation.sh remind --request-id setup-002 --actor-scope telegram:setup
4) Apply Telegram policy defaults from config/openclaw/telegram.example.jsonc:
   - dmPolicy=pairing
   - groups require mention
   - allowActions=[status,remind]
   - denyUnknownActions=true
5) Register skill from config/openclaw/SKILL.md with intent mapping:
   - status -> scripts/openclaw/openppl-automation.sh status
   - send reminder -> scripts/openclaw/openppl-automation.sh remind --request-id <request_id> --actor-scope <actor_scope>
6) (Optional) Configure MCP route from config/openclaw/mcp-server.example.jsonc using wrapper script as command target.
7) Run: bash tests/openclaw/smoke_test.sh

Definition of done:
- dry-run passes
- wrapper status/remind calls pass
- smoke test passes
- skill is active with only status/remind intents
- unknown action requests are denied

Output format:
- Show a short checklist with pass/fail per step.
- Show exact commands run.
- If any failure occurs, provide minimal fix and rerun only failed step.
```

## 1) Validate Local Integration Assets

Run the built-in dry-run check:

```bash
bash scripts/openclaw/deploy.sh dry-run
```

Expected output:

```text
dry-run checks passed
```

## 2) Verify openppl Automation Commands

Check status path:

```bash
openppl automation status
```

Check reminder action path:

```bash
openppl automation action --name remind --request-id test-001 --actor-scope telegram:manual
```

Run same reminder again with same request ID to confirm idempotency (`result_state` should become `replayed`).

## 3) Use the OpenClaw Wrapper (Recommended)

Wrapper enforces allowlisted subcommands and argument validation.

Status:

```bash
bash scripts/openclaw/openppl-automation.sh status
```

Reminder:

```bash
bash scripts/openclaw/openppl-automation.sh remind --request-id tg-001 --actor-scope telegram:default
```

## 4) Configure Telegram Channel in OpenClaw

Use `config/openclaw/telegram.example.jsonc` as your base and copy values into your OpenClaw runtime config.

Key defaults:

- `dmPolicy: "pairing"`
- `groups.*.requireMention: true`
- `policy.allowActions: ["status", "remind"]`
- `policy.denyUnknownActions: true`

Set your token in environment (do not commit real tokens):

```bash
export OPENCLAW_TELEGRAM_BOT_TOKEN="<your-bot-token>"
```

If your OpenClaw flow requires pairing approval, approve the code in OpenClaw:

```bash
openclaw pairing approve telegram <CODE>
```

## 5) Add the Skill to OpenClaw

Use `config/openclaw/SKILL.md` as the source for your OpenClaw skill definition.

Intent mapping in this repo:

- `status` -> `scripts/openclaw/openppl-automation.sh status`
- `send reminder` -> `scripts/openclaw/openppl-automation.sh remind --request-id <request_id> --actor-scope <actor_scope>`

Important constraints:

- Unknown intents must be denied
- No raw shell passthrough
- No direct free-form `openppl` command execution

## 6) Optional MCP Server Route

If you want OpenClaw to call `openppl` through MCP stdio transport, start from:

- `config/openclaw/mcp-server.example.jsonc`

This points MCP calls to:

```text
bash scripts/openclaw/openppl-automation.sh
```

Keep the same allowlist (`status`, `remind`) and deny-unknown policy.

## 7) Run End-to-End Smoke Tests

```bash
bash tests/openclaw/smoke_test.sh
```

Expected output:

```text
openclaw smoke test passed
```

## Telegram Usage Examples

After skill/config is active in OpenClaw, test from Telegram:

- "status"
- "send reminder"

If your bot is in groups, mention-gating should require tagging the bot before actioning requests.

## Troubleshooting

- `wrapper.invalid_arguments`: verify wrapper args (`status` takes no extra args, `remind` requires `--request-id`)
- `openppl: command not found`: set `OPENPPL_BIN` to your binary path
- unexpected denied action: confirm OpenClaw policy allowlist includes `status` and `remind`
- duplicate reminder behavior not seen: ensure same `--request-id` is reused for retry tests

## Security Notes

- Never commit Telegram bot tokens or runtime secrets
- Keep pairing enabled for DMs and mention requirement enabled for groups
- Keep deny-by-default behavior for unknown actions
