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
