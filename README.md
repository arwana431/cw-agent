# CertWatch Agent

[![CI](https://github.com/certwatch-app/cw-agent/actions/workflows/ci.yml/badge.svg)](https://github.com/certwatch-app/cw-agent/actions/workflows/ci.yml)
[![Release](https://github.com/certwatch-app/cw-agent/actions/workflows/release.yml/badge.svg)](https://github.com/certwatch-app/cw-agent/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/certwatch-app/cw-agent)](https://goreportcard.com/report/github.com/certwatch-app/cw-agent)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

SSL/TLS certificate monitoring agent for [CertWatch](https://certwatch.app). Monitor certificates on your infrastructure and sync data to the CertWatch cloud platform.

## Features

- **Interactive setup wizard** - `cw-agent init` guides you through configuration
- **Config-driven monitoring** - Define certificates to monitor in a YAML file
- **Automatic scanning** - Continuously scan certificates at configurable intervals
- **Cloud sync** - Automatically sync certificate data to CertWatch dashboard
- **Agent state persistence** - Agent ID survives restarts, supports name changes
- **Smart certificate migration** - Certificates transfer when resetting agents
- **Chain validation** - Detect chain issues, expiration, and weak cryptography
- **Lightweight** - Single binary, minimal resource usage
- **Secure** - Runs without root, distroless Docker image
- **Beautiful CLI** - Polished terminal UI with colors and styling

## Quick Start

### 1. Install

**Quick install (Linux/macOS):**
```bash
curl -sSL https://certwatch.app/install.sh | bash
```

**Using Go:**
```bash
go install github.com/certwatch-app/cw-agent/cmd/cw-agent@latest
```

**Using Homebrew (macOS/Linux):**
```bash
brew install certwatch-app/tap/cw-agent
```

**Using Docker:**
```bash
docker pull ghcr.io/certwatch-app/cw-agent:latest
```

**Download binary:**

Download the latest release from the [releases page](https://github.com/certwatch-app/cw-agent/releases).

### 2. Configure

**Interactive wizard (recommended):**
```bash
cw-agent init
```

This launches an interactive wizard that guides you through:
- Setting the config file path
- Entering your CertWatch API key
- Configuring agent name and intervals
- Adding certificates to monitor

**Or manually** create a `certwatch.yaml` configuration file:

```yaml
api:
  endpoint: "https://api.certwatch.app"
  key: "cw_xxxxxxxx_xxxx..."  # Get from CertWatch dashboard

agent:
  name: "production-monitor"
  sync_interval: 5m
  scan_interval: 1m

certificates:
  - hostname: "www.example.com"
    port: 443
    tags: ["production", "web"]

  - hostname: "api.example.com"
    port: 443
    tags: ["production", "api"]
```

See [certwatch.example.yaml](certwatch.example.yaml) for a complete example.

### 3. Run

```bash
cw-agent start -c certwatch.yaml
```

Or with Docker:

```bash
docker run -v $(pwd)/certwatch.yaml:/etc/certwatch/certwatch.yaml \
  ghcr.io/certwatch-app/cw-agent:latest
```

## Commands

### `cw-agent init`

Interactive configuration wizard:

```bash
# Launch interactive wizard
cw-agent init

# Specify output path
cw-agent init -o /etc/certwatch/certwatch.yaml

# Non-interactive mode (for CI/automation)
CW_API_KEY=cw_xxx CW_CERTIFICATES=example.com,api.example.com \
  cw-agent init --non-interactive -o certwatch.yaml
```

### `cw-agent start`

Start the monitoring agent:

```bash
# Start with config file
cw-agent start -c certwatch.yaml

# Reset agent state (creates new agent, migrates certificates)
cw-agent start -c certwatch.yaml --reset-agent

# Skip confirmation prompts (for CI/automation)
cw-agent start -c certwatch.yaml --reset-agent --yes
```

### `cw-agent validate`

Validate configuration without starting:

```bash
cw-agent validate -c certwatch.yaml
```

### `cw-agent version`

Show version information:

```bash
cw-agent version
```

## Agent State & Migration

The agent persists its state in a `.certwatch-state.json` file alongside your config file. This enables:

- **Restart resilience** - Agent ID survives restarts
- **Name change detection** - Warns if you change `agent.name` in config
- **Certificate migration** - When resetting, certificates transfer to new agent

### Changing Agent Name

If you change `agent.name` in your config, the agent will warn you:

```
! Agent name changed: "old-name" → "new-name"

Options:
  1. Continue with new name (migrates matching certs):
     cw-agent start -c certwatch.yaml --reset-agent

  2. Keep existing agent (revert config):
     Edit certwatch.yaml: agent.name = "old-name"
```

Using `--reset-agent`:
- Creates a new agent with the new name
- Migrates certificates that are still in your config to the new agent
- Orphans certificates no longer in your config (they become dashboard-managed)

## Configuration

### API Settings

| Field | Description | Default |
|-------|-------------|---------|
| `api.endpoint` | CertWatch API URL | `https://api.certwatch.app` |
| `api.key` | API key with `cloud:sync` scope | Required |
| `api.timeout` | HTTP request timeout | `30s` |

### Agent Settings

| Field | Description | Default |
|-------|-------------|---------|
| `agent.name` | Unique name for this agent | `default-agent` |
| `agent.sync_interval` | How often to sync with cloud | `5m` |
| `agent.scan_interval` | How often to scan certificates | `1m` |
| `agent.concurrency` | Max concurrent scans | `10` |
| `agent.log_level` | Log level (debug/info/warn/error) | `info` |

### Certificate Settings

| Field | Description | Default |
|-------|-------------|---------|
| `hostname` | Hostname to connect to | Required |
| `port` | Port to connect to | `443` |
| `tags` | Tags for organization | `[]` |
| `notes` | Notes about this certificate | `""` |

## Environment Variables

Configuration can also be set via environment variables with the `CW_` prefix:

```bash
export CW_API_KEY="cw_xxxx..."
export CW_AGENT_NAME="my-agent"
```

For non-interactive init:

| Variable | Required | Description |
|----------|----------|-------------|
| `CW_API_KEY` | Yes | API key with `cloud:sync` scope |
| `CW_API_ENDPOINT` | No | API endpoint URL |
| `CW_AGENT_NAME` | No | Agent name |
| `CW_SYNC_INTERVAL` | No | Sync interval (e.g., `5m`) |
| `CW_SCAN_INTERVAL` | No | Scan interval (e.g., `1m`) |
| `CW_LOG_LEVEL` | No | Log level |
| `CW_CERTIFICATES` | Yes | Comma-separated hostnames |

## Running as a Service

### systemd (Linux)

Create `/etc/systemd/system/cw-agent.service`:

```ini
[Unit]
Description=CertWatch Agent
After=network.target

[Service]
Type=simple
User=certwatch
ExecStart=/usr/local/bin/cw-agent start -c /etc/certwatch/certwatch.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable cw-agent
sudo systemctl start cw-agent
```

### Docker Compose

```yaml
version: '3.8'
services:
  cw-agent:
    image: ghcr.io/certwatch-app/cw-agent:latest
    restart: unless-stopped
    volumes:
      - ./certwatch.yaml:/etc/certwatch/certwatch.yaml:ro
```

## Getting an API Key

1. Log in to [CertWatch](https://certwatch.app)
2. Go to **Settings** > **API Keys**
3. Create a new key with the `cloud:sync` scope
4. Copy the key (it's only shown once!)

## Development

### Building

```bash
# Build
make build

# Run tests
make test

# Run linter
make lint

# Build for all platforms
make build-all
```

### Project Structure

```
cw-agent/
├── cmd/cw-agent/       # Entry point
├── internal/
│   ├── agent/          # Main orchestrator
│   ├── cmd/            # CLI commands
│   │   └── initcmd/    # Init wizard forms & logic
│   ├── config/         # Configuration loading
│   ├── scanner/        # TLS certificate scanning
│   ├── state/          # Agent state persistence
│   ├── sync/           # API client
│   ├── ui/             # Shared CLI styling
│   └── version/        # Version info
├── certwatch.example.yaml
├── Dockerfile
├── Makefile
└── README.md
```

## Changelog

### v0.2.1 (Current)

- **Agent state persistence** - Agent ID stored in `.certwatch-state.json`
- **Name change detection** - Warns when `agent.name` changes in config
- **`--reset-agent` flag** - Reset state and migrate certificates to new agent
- **`--yes` flag** - Skip confirmation prompts for CI/automation
- **Unified CLI styling** - All commands now have consistent, polished output
- **Smart certificate migration** - Certificates transfer during agent reset

### v0.2.0

- **`cw-agent init` command** - Interactive configuration wizard
- **Non-interactive mode** - Create configs via environment variables
- **Beautiful forms** - Powered by charmbracelet/huh

### v0.1.0

- Initial release
- Certificate scanning and cloud sync
- `start`, `validate`, `version` commands
- Docker and systemd support

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

Apache 2.0 - see [LICENSE](LICENSE) for details.
