# openppl

The open source PPL study planner for terminal, web, and automation.

`openppl` helps you plan and track training with:

- interactive TUI workflows
- web dashboard mode
- non-interactive automation commands for integrations (OpenClaw/Telegram)

---

## Installation

### macOS / Linux

```bash
# Recommended installer
curl -fsSL https://app.openppl.net/install | bash
```

Alternative:

```bash
# Build a local binary from source
git clone git@github.com:toomij99/openppl.git
cd openppl
go build -o openppl .

# Run without installing
go run .
```

### Windows (PowerShell)

Recommended: use WSL and run the macOS/Linux installer command above.

Native source build:

```powershell
git clone https://github.com/toomij99/openppl.git
cd openppl
go build -o openppl.exe .
.\openppl.exe help
```

### Any Platform (Go toolchain)

```bash
# Install with Go (requires Go 1.25+)
go install ./...
```

After install:

```bash
openppl help
```

---

## Quick Start

```bash
# Launch TUI (first run triggers onboarding)
openppl

# Run onboarding manually
openppl onboard

# Reconfigure core settings
openppl --configure
```

---

## Commands

```bash
# Show help
openppl help

# Show recent runtime errors
openppl logs

# Launch web mode
openppl web

# Launch web mode with custom bind settings
openppl web --hostname 0.0.0.0 --port 5016

# Automation status (JSON output)
openppl automation status

# Automation action (idempotent reminder)
openppl automation action --name remind --request-id req-001 --actor-scope telegram:default
```

---

## OpenClaw Integration

For Telegram integration, skill setup, policy defaults, MCP example config, and smoke/deploy checks:

- `docs/README_openclaw.md`

---

## Environment

Create local env file from template:

```bash
cp .env.example .env
```

Set your own values in `.env` and never commit secrets.

---

## Development

```bash
go test ./...
go build ./...
```
