# Supabase Management Tools

A comprehensive toolkit for managing self-hosted Supabase instances, consisting of two complementary tools:

1. **supactl** - A modern CLI for managing Supabase instances via a SupaControl server
2. **supascale.sh** - A bash script for direct local management of multiple Supabase instances

[![Test](https://github.com/yourusername/supactl/workflows/Test/badge.svg)](https://github.com/yourusername/supactl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/supactl)](https://goreportcard.com/report/github.com/yourusername/supactl)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## 🚀 Quick Start

### supactl (Recommended)

The modern way to manage Supabase instances through a central server.

```bash
# Install
curl -L https://github.com/yourusername/supactl/releases/latest/download/supactl-linux-amd64 -o supactl
chmod +x supactl
sudo mv supactl /usr/local/bin/

# Login and start using
supactl login https://your-supacontrol-server.com
supactl create my-project
supactl list
```

**[Full supactl Documentation →](SUPACTL_README.md)**

### supascale.sh

Direct local management for single-machine deployments.

```bash
# Clone and use
git clone https://github.com/yourusername/supactl.git
cd supactl
chmod +x supascale.sh
./supascale.sh add my-project
```

## 📋 Table of Contents

- [Overview](#overview)
- [Tools Comparison](#tools-comparison)
- [Installation](#installation)
- [Quick Examples](#quick-examples)
- [Features](#features)
- [Documentation](#documentation)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This repository provides two powerful tools for managing self-hosted Supabase instances:

### supactl - Modern CLI Tool

`supactl` is a Go-based command-line interface that connects to a SupaControl management server, providing:

- **Centralized Management**: Manage instances across multiple servers from one place
- **API-Driven**: Full REST API integration with the SupaControl server
- **Local Linking**: Connect local development directories to remote instances
- **Cross-Platform**: Single binary for Linux, macOS, and Windows
- **Interactive**: Beautiful CLI with interactive prompts
- **Secure**: Encrypted credentials storage with proper permissions

### supascale.sh - Direct Management Script

`supascale.sh` is a bash script for direct local management of Supabase instances:

- **Single Machine Focus**: Perfect for managing multiple instances on one server
- **Docker Integration**: Direct docker-compose management
- **Port Management**: Intelligent automatic port allocation
- **No Server Required**: Standalone script, no additional infrastructure
- **Quick Setup**: Automated Supabase project initialization

## 🔄 Tools Comparison

| Feature | supactl | supascale.sh |
|---------|---------|--------------|
| **Architecture** | Client-Server | Standalone |
| **Best For** | Multi-server, team usage | Single machine, local dev |
| **Setup Complexity** | Requires SupaControl server | Just the script |
| **Remote Management** | ✅ Yes | ❌ No |
| **API Integration** | ✅ Full REST API | ❌ N/A |
| **Docker Required** | On server only | ✅ Local required |
| **Cross-Platform** | ✅ Yes | Linux/macOS only |
| **Interactive CLI** | ✅ Beautiful prompts | Basic bash |
| **Local Development** | ✅ Link feature | ✅ Direct access |

**Choose supactl if:** You're managing multiple servers, working in a team, or want centralized control.

**Choose supascale.sh if:** You're managing instances on a single machine and prefer direct Docker access.

## 📦 Installation

### Installing supactl

#### Using curl (Linux/macOS)
```bash
# Linux
curl -L https://github.com/yourusername/supactl/releases/latest/download/supactl-linux-amd64 -o supactl
chmod +x supactl
sudo mv supactl /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yourusername/supactl/releases/latest/download/supactl-darwin-amd64 -o supactl
chmod +x supactl
sudo mv supactl /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/yourusername/supactl/releases/latest/download/supactl-darwin-arm64 -o supactl
chmod +x supactl
sudo mv supactl /usr/local/bin/
```

#### Using Homebrew (macOS)
```bash
brew tap yourusername/tap
brew install supactl
```

#### From Source
```bash
git clone https://github.com/yourusername/supactl.git
cd supactl
make build
sudo make install
```

### Installing supascale.sh

```bash
git clone https://github.com/yourusername/supactl.git
cd supactl
chmod +x supascale.sh

# Optional: Add to PATH
sudo ln -s $(pwd)/supascale.sh /usr/local/bin/supascale
```

**Prerequisites for supascale.sh:**
- Docker and Docker Compose
- `jq` (JSON processor)
- Bash shell environment
- Supabase CLI

## 🎨 Quick Examples

### Using supactl

```bash
# One-time login
supactl login https://supacontrol.example.com

# Create and manage instances
supactl create production
supactl create staging
supactl list

# Link to local directory
cd ~/my-project
supactl link
supactl status

# When done
supactl delete staging
```

### Using supascale.sh

```bash
# Create instances
./supascale.sh add production
./supascale.sh add staging

# Manage lifecycle
./supascale.sh start production
./supascale.sh list
./supascale.sh stop production

# Cleanup
./supascale.sh remove staging
```

## ✨ Features

### supactl Features

- **Authentication & Configuration**: Secure login with API key authentication
- **Instance Management**: Create, list, delete Supabase instances
- **Local Project Linking**: Link development directories to remote instances
- **Status Monitoring**: View detailed instance information
- **Security**: Credentials stored with 600 permissions, no key echoing
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Single Binary**: No runtime dependencies

### supascale.sh Features

- **Easy Project Creation**: Automated setup with unique configurations
- **Secure by Default**: Auto-generated passwords and secrets
- **Port Management**: Intelligent automatic port allocation (base + 1000 per project)
- **Container Management**: Start, stop, and manage Docker containers
- **Configuration Management**: Centralized JSON-based storage
- **Docker Integration**: Seamless docker-compose integration
- **Update Mechanism**: Self-updating capability

## 📚 Documentation

### Detailed Documentation

- **[supactl Full Documentation](SUPACTL_README.md)** - Complete guide for the CLI tool
  - All commands and options
  - API endpoints reference
  - Configuration details
  - Troubleshooting guide

- **[supascale.sh Guide](#supascalesh-detailed-guide)** - See below for bash script documentation

### API Documentation

supactl works with SupaControl servers that implement these endpoints:

- `GET /api/v1/auth/me` - Validate authentication
- `GET /api/v1/instances` - List instances
- `POST /api/v1/instances` - Create instance
- `GET /api/v1/instances/<name>` - Get instance details
- `DELETE /api/v1/instances/<name>` - Delete instance

## 🛠 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/supactl.git
cd supactl

# Build
make build

# Run tests
make test

# Build for all platforms
make build-all

# Install locally
make install
```

### Running Tests

```bash
# Run all tests
go test ./...

# With coverage
make test-coverage

# Run specific tests
go test -v ./internal/api
```

### Project Structure

```
.
├── cmd/                 # Cobra command definitions
│   ├── create.go
│   ├── delete.go
│   ├── link.go
│   ├── list.go
│   ├── login.go
│   ├── logout.go
│   ├── status.go
│   └── unlink.go
├── internal/
│   ├── api/            # API client implementation
│   ├── auth/           # Authentication & config
│   └── link/           # Local project linking
├── .github/
│   └── workflows/      # CI/CD pipelines
├── main.go
├── Makefile
├── supascale.sh        # Bash management script
└── README.md
```

## supascale.sh Detailed Guide

### Prerequisites

- Docker and Docker Compose
- `jq` (JSON processor)
- Git
- Bash shell environment
- Sudo privileges (required for Docker operations)
- Supabase CLI (must be in your PATH)

### Available Commands

```bash
./supascale.sh [command] [options]
```

- `list` - Display all configured projects
- `add <project_id>` - Create a new Supabase instance
- `start <project_id>` - Start a specific project
- `stop <project_id>` - Stop a specific project
- `remove <project_id>` - Remove a project from configuration
- `help` - Show help message

### Configuration

The script uses a central configuration file at `$HOME/.supascale_database.json`:

```json
{
  "projects": {
    "project-id": {
      "directory": "/path/to/project",
      "ports": {
        "api": 54321,
        "db": 54322,
        "shadow": 54320,
        "studio": 54323,
        "inbucket": 54324,
        "smtp": 54325,
        "pop3": 54326,
        "analytics": 54327,
        "pooler": 54329,
        "kong_https": 54764
      }
    }
  },
  "last_port_assigned": 54321
}
```

### Port Allocation

Base port starts at 54321 and increments by 1000 per project:

- Shadow Port: Base - 1
- API Port (Kong): Base
- Database Port: Base + 1
- Studio Port: Base + 2
- Inbucket Port: Base + 3
- SMTP Port: Base + 4
- POP3 Port: Base + 5
- Analytics Port: Base + 6
- Pooler Port: Base + 8
- Kong HTTPS Port: Base + 443

### Security

- Passwords generated using `/dev/urandom` (40 chars, alphanumeric)
- Credentials stored in project-specific `.env` files
- Each project isolated in Docker containers
- JWT secrets auto-generated for Supabase

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Install dependencies
make deps

# Run linter
make lint

# Run tests
make test

# Build
make build
```

## 📄 License

GPL V3 License - Copyright (c) 2025 Frog Byte, LLC

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

See [LICENSE](LICENSE) for full details.

## 🙏 Acknowledgments

- Original development by Frog Byte, LLC
- Built on top of [Supabase](https://supabase.com/)
- CLI framework powered by [Cobra](https://github.com/spf13/cobra) 
