# supactl

A unified, context-aware CLI for managing self-hosted Supabase instances. Supports both remote management via SupaControl servers and local Docker-based management in a single binary.

[![Test](https://img.shields.io/github/actions/workflow/status/qubitquilt/supactl/test.yml?style=for-the-badge&label=tests)](https://github.com/qubitquilt/supactl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/qubitquilt/supactl?style=for-the-badge)](https://goreportcard.com/report/github.com/qubitquilt/supactl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/qubitquilt/supactl?style=for-the-badge)](https://github.com/qubitquilt/supactl/releases)

## Overview

`supactl` provides a kubectl-inspired interface for Supabase instance management. Switch seamlessly between local (Docker) and remote (SupaControl server) modes using contexts. All core commands work identically across providers.

### Key Features
- **Context-Aware**: Unified commands for local and remote management
- **kubectl-Style UX**: Familiar `get`, `describe`, `config` subcommands
- **Local Mode**: Direct Docker Compose management, no server required
- **Remote Mode**: API-driven control of SupaControl servers
- **Project Linking**: Link local directories to remote instances
- **Cross-Platform**: Linux, macOS, Windows
- **Secure**: Encrypted secrets, no credential echoing
- **Single Binary**: ~10MB, no runtime dependencies

## Installation

### Recommended: Installation Script
```bash
curl -sSL https://raw.githubusercontent.com/qubitquilt/supactl/main/scripts/install.sh | bash
```
- Auto-detects platform/architecture
- Verifies checksums
- Installs to `/usr/local/bin`
- Supports Linux (amd64/arm64), macOS (Intel/Apple Silicon)

### Homebrew (macOS)
```bash
brew tap qubitquilt/homebrew-tap
brew install supactl
```

### Manual Download
Download from [Releases](https://github.com/qubitquilt/supactl/releases/latest):
- Linux amd64: `supactl_Linux_x86_64.tar.gz`
- Linux arm64: `supactl_Linux_arm64.tar.gz`
- macOS Intel: `supactl_Darwin_x86_64.tar.gz`
- macOS Apple Silicon: `supactl_Darwin_arm64.tar.gz`
- Windows: `supactl_Windows_x86_64.zip` or `supactl_Windows_arm64.zip`

Extract, `chmod +x supactl`, and add to PATH.

### From Source
```bash
git clone https://github.com/qubitquilt/supactl.git
cd supactl
make build
sudo make install  # or copy to PATH
```

## Quick Start

### 1. Setup Contexts
```bash
# Local Docker context (default)
supactl config set-context local --provider=local
supactl config use-context local

# Remote SupaControl context
supactl config set-context production --provider=remote --server=https://your-supacontrol.com --api-key=sk_...
supactl config use-context production
```

Or use `login` for quick remote setup:
```bash
supactl login https://your-supacontrol.com
# Enter API key when prompted
```

### 2. Manage Instances
Commands work in the current context:
```bash
# Create instance
supactl create my-project

# List instances
supactl list
# or kubectl-style
supactl get instances

# Start/Stop
supactl start my-project
supactl stop my-project

# Details
supactl describe instance my-project
supactl logs my-project --lines=50
```

### 3. Local Mode Example
```bash
supactl config use-context local
supactl local add my-local-project  # Clones Supabase, sets up Docker
supactl local start my-local-project
# Access Studio at http://localhost:54323 (auto-assigned port)
```

### 4. Linking (Remote Mode)
```bash
cd ~/my-project
supactl link  # Select from instances
supactl status  # View linked details
```

## Configuration

Contexts are stored in `~/.supacontrol/config.json` (0600 permissions).

Example:
```json
{
  "current-context": "production",
  "contexts": {
    "local": { "provider": "local" },
    "production": { "provider": "remote", "server_url": "https://...", "api_key": "sk_..." }
  }
}
```

Legacy v1 format auto-migrates.

### Context Commands
```bash
supactl config get-contexts     # List contexts
supactl config use-context <name>  # Switch
supactl config current-context   # Show active
supactl config set-context <name> --provider=local|remote [options]  # Create/update
supactl config delete-context <name>  # Remove
```

## Commands

### Core Instance Management (Context-Aware)
These work in local or remote contexts, but `create` is remote-only (use `local add` for local creation):

- `supactl create <name>`: Create new remote instance
  - Remote: Calls API to provision. Local: Not supported; use `supactl local add <name>` instead.
  - Name regex: `^[a-z0-9][a-z0-9_-]*$`

- `supactl list`: List instances (tabular)
- `supactl delete <name>`: Delete instance (confirmation prompt)
- `supactl start <name>`: Start instance
- `supactl stop <name>`: Stop instance
- `supactl restart <name>`: Restart instance
- `supactl logs <name> [--lines=N]`: View recent logs

### kubectl-Style Commands
- `supactl get instances`: List in table format (alias: `list`)
- `supactl describe instance <name>`: Detailed info (status, URLs, ports, etc.)

### Local Subcommands
Dedicated local management (ignores remote context):
- `supactl local add <name>`: Create local project
- `supactl local list`: List local projects/ports
- `supactl local start <name>`: Start Docker services
- `supactl local stop <name>`: Stop services
- `supactl local remove <name>`: Remove from database (keeps files)

### Linking & Status (Remote Mode)
- `supactl link`: Link current dir to instance (creates `.supacontrol/project`)
- `supactl unlink`: Remove link
- `supactl status`: Show linked instance details (URLs, keys, etc.)

### Authentication
- `supactl login <server_url>`: Setup default remote context, prompt for API key
- `supactl logout`: Clear credentials

### Help & Version
- `supactl --help`: All commands
- `supactl <cmd> --help`: Specific help
- `supactl --version`: Show version

## Examples

### Full Local Workflow
```bash
supactl config use-context local
supactl local add dev-app
supactl get instances
supactl start dev-app
supactl describe instance dev-app  # Shows ports: API=54321, Studio=54323, etc.
# Develop at http://localhost:54321
supactl stop dev-app
supactl delete dev-app
```

### Remote Workflow with Linking
```bash
supactl login https://supacontrol.example.com
supactl create prod-app
cd ~/prod-app
supactl link  # Select prod-app
supactl status  # Shows remote URLs, anon key, etc.
supactl get instances  # From any dir, uses link if present
```

### Multi-Context
```bash
supactl config set-context dev --provider=remote --server=https://dev.example.com --api-key=sk_dev
supactl config set-context staging --provider=remote --server=https://staging.example.com --api-key=sk_staging
supactl config use-context dev
supactl create dev-instance
supactl config use-context staging
supactl create staging-instance  # Independent
```

## Local Mode Details

- **Storage**: `~/.supascale_database.json` (0600 perms)
- **Ports**: Auto-allocated (base 54321 + 1000 * project_index)
  - API: base, DB: base+1, Studio: base+2, etc.
- **Secrets**: Auto-generated (crypto/rand, HS256 JWT)
- **Isolation**: Per-project Docker networks/containers
- **Supersedes**: Legacy `supascale.sh` (compatible DB format)

## Remote Mode Details

- **Linking**: `.supacontrol/project` file (git-ignored)
- **Auth**: Bearer token (API key), HTTPS required
- **Validation**: `GET /api/v1/auth/me` on login

## API Endpoints (SupaControl Server)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Validate API key |
| GET | `/api/v1/instances` | List instances |
| POST | `/api/v1/instances` | Create instance |
| GET | `/api/v1/instances/{name}` | Get details |
| DELETE | `/api/v1/instances/{name}` | Delete |
| POST | `/api/v1/instances/{name}/start` | Start |
| POST | `/api/v1/instances/{name}/stop` | Stop |
| POST | `/api/v1/instances/{name}/restart` | Restart |
| GET | `/api/v1/instances/{name}/logs?lines=N` | Get logs |

All use `Authorization: Bearer <api_key>`.

## Troubleshooting

### Common Errors
- **"Not logged in"**: Run `login` or `config set-context` with credentials.
- **"Invalid context"**: Use `config get-contexts`; switch with `use-context`.
- **"No project linked"** (status/link): Run `link` in project dir.
- **Port conflicts** (local): Check `docker ps`; use `local remove` for cleanup.
- **Auth failed**: Verify API key/server URL; test connectivity.
- **Invalid name**: Use lowercase alphanum + hyphens, start/end alphanumeric.

### Logs & Debug
- `supactl logs <name> --lines=100`
- Set `SUPACTL_DEBUG=1` for verbose output.
- Config location: `~/.supacontrol/config.json` (check perms).

### Windows Notes
- Paths use `%USERPROFILE%\.supacontrol\`
- Docker requires WSL2 or Hyper-V.
- No file perm enforcement (uses ACLs).

## Development

### Build & Test
```bash
make build      # Current platform
make build-all  # All platforms
make test       # Unit tests
make lint       # golangci-lint
make test-coverage
```

### Project Structure
```
.
├── cmd/          # Cobra commands
│   ├── root.go
│   ├── config.go, create.go, delete.go, ...
│   └── local_*.go
├── internal/
│   ├── api/      # Remote API client
│   ├── auth/     # Config/auth
│   ├── link/     # Project linking
│   ├── local/    # Docker/local mgmt
│   └── provider/ # Abstraction layer
├── scripts/      # install.sh, uninstall.sh
├── main.go
├── Makefile
├── go.mod
└── ...
```

### Dependencies
- Cobra (CLI)
- Survey (prompts)
- Standard lib (HTTP, JSON, crypto)

### Testing
- `go test ./...`
- Cross-platform: `GOOS=windows go test ./...`
- Coverage: `make test-coverage`

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md).

1. Fork & branch
2. Implement + tests
3. Lint & build
4. PR with description

## License
MIT - © 2025 Qubit Quilt. See [LICENSE](LICENSE).

## Acknowledgments
- [Supabase](https://supabase.com/)
- [Cobra](https://github.com/spf13/cobra)
- Community contributions

## Support
- Issues: [GitHub](https://github.com/qubitquilt/supactl/issues)
- Docs: This file + [README.md](README.md)

---
Made with ❤️ by Qubit Quilt
