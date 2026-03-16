# openppl

The open source PPL study planner for terminal, web, and automation.

Prepare smarter, fly more confidently. `openppl` helps pilot students turn checkride prep into a clear weekly plan with realistic training milestones, budget visibility, and consistent progress tracking.

Built for PPL today, with a roadmap to support your full training journey: CPL, IFR, ME, CFI, and potentially ATPL.

`openppl` helps you plan and track training with:

- interactive TUI workflows
- web dashboard mode
- non-interactive automation commands for integrations (OpenClaw/Telegram)

Use it to:

- break down your checkride preparation into manageable study and flight tasks
- track expected training costs and avoid budget surprises
- keep momentum from first lesson to checkride day

---

## Installation

### macOS / Linux

```bash
# Recommended installer
curl -fsSL https://openppl.happycloud.ru/install | bash
```

The installer shows an interactive summary with ASCII logo, installed/new version, and a quick command guide before downloading.

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

```powershell
# Recommended installer
powershell -NoProfile -ExecutionPolicy Bypass -Command "iwr -useb https://openppl.happycloud.ru/install.ps1 | iex"
```

Pin a version:

```powershell
$env:OPENPPL_VERSION = "v0.1.11"
powershell -NoProfile -ExecutionPolicy Bypass -Command "iwr -useb https://openppl.happycloud.ru/install.ps1 | iex"
```

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

# Show MOTD ACS daily quiz card
openppl motd

# Run today's multiple-choice quiz now
openppl motd quiz

# Show checkride readiness progress
openppl motd progress

# Show weakest ACS areas
openppl motd weak
```

---

## MOTD Daily Quiz (Ubuntu)

`openppl` can show an ACS-focused quiz during Ubuntu login and track your progress toward PPL checkride readiness.

Install MOTD integration (root required):

```bash
sudo openppl motd install
```

What it installs:

- `/etc/update-motd.d/99-openppl-acs` for login-time daily ACS card
- `/etc/profile.d/openppl-recall.sh` for interactive daily quiz prompt

Manual commands:

```bash
# Show today's quiz card
openppl motd

# Answer today's quiz question
openppl motd quiz

# View readiness score, accuracy, and area breakdown
openppl motd progress

# View lowest-performing ACS areas first
openppl motd weak
```

Tip: if a newer release is available, `openppl motd` shows an update recommendation with the exact installer command.

Disable login quiz prompt for a shell session:

```bash
export OPENPPL_MOTD_RECALL=0
```

---

## OpenClaw Integration

For Telegram integration, skill setup, policy defaults, MCP example config, and smoke/deploy checks:

- `docs/README_openclaw.md`

---

## Deploy on Ubuntu (DigitalOcean) with Caddy + SSL

Example target:

- Ubuntu 24.04 droplet
- domain: `app.example.com`
- app runs on `127.0.0.1:5016`
- Caddy terminates TLS and reverse proxies to openppl web mode

### 1) Server prep

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl ca-certificates ufw caddy

# Firewall
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable
```

### 2) Install openppl

```bash
curl -fsSL https://openppl.happycloud.ru/install | bash
openppl help
```

If needed, ensure binary is in your shell path (for example `~/bin`).

### 3) First-run setup (once)

```bash
openppl onboard
```

### 4) Create systemd service

Create `/etc/systemd/system/openppl.service`:

```ini
[Unit]
Description=openppl web service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root
ExecStart=/root/bin/openppl web --hostname 127.0.0.1 --port 5016
Restart=always
RestartSec=5
Environment=HOME=/root

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable openppl
sudo systemctl start openppl
sudo systemctl status openppl
```

### 5) Configure Caddy for HTTPS

Create `/etc/caddy/Caddyfile`:

```caddy
app.example.com {
    reverse_proxy 127.0.0.1:5016
}
```

Reload Caddy:

```bash
sudo caddy validate --config /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

Caddy will automatically provision and renew TLS certificates when DNS for `app.example.com` points to your droplet.

### 6) Operations

```bash
# App logs
sudo journalctl -u openppl -f

# Caddy logs
sudo journalctl -u caddy -f

# Restart services
sudo systemctl restart openppl
sudo systemctl restart caddy
```

### 7) Optional hardening

- Run service as a dedicated non-root user
- Store env in an `EnvironmentFile=` and lock file permissions
- Keep app bound to `127.0.0.1` only (as shown above)
- Add DigitalOcean cloud firewall rules for `22`, `80`, `443` only

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
