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
curl -fsSL https://openppl.happycloud.ru/install | bash
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
