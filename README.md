# Supabase Management Tools

A comprehensive toolkit for managing self-hosted Supabase instances.

**supactl** is a modern, unified CLI for managing Supabase instances both remotely (via a SupaControl server) and locally (direct Docker management).

[![Test](https://img.shields.io/github/actions/workflow/status/qubitquilt/supactl/test.yml?style=for-the-badge&label=tests)](https://github.com/qubitquilt/supactl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/qubitquilt/supactl?style=for-the-badge)](https://goreportcard.com/report/github.com/qubitquilt/supactl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/qubitquilt/supactl?style=for-the-badge)](https://github.com/qubitquilt/supactl/releases)

## Built With

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Go Version](https://img.shields.io/github/go-mod/go-version/qubitquilt/supactl?style=for-the-badge)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![macOS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0)
![Windows](https://img.shields.io/badge/Windows-0078D4?style=for-the-badge&logo=windows&logoColor=white)

## 🚀 Quick Start

### Install supactl

```bash
# Install via script (recommended - secure, automated, cross-platform)
curl -sSL https://raw.githubusercontent.com/qubitquilt/supactl/main/scripts/install.sh | bash
```

### Context-Aware Management

supactl uses kubectl-style contexts to switch between remote and local management:

```bash
# Set up local context (Docker-based)
supactl config set-context local --provider=local
supactl config use-context local

# Set up remote context (SupaControl server)
supactl config set-context production --provider=remote --server=https://your-supacontrol-server.com --api-key=sk_...
supactl config use-context production
```

### Basic Usage

```bash
# List all contexts
supactl config get-contexts

# Create and manage instances (create only in remote; use 'local add' for local)
supactl local add my-project
supactl list
supactl start my-project

# kubectl-style commands
supactl get instances
supactl describe instance my-project
```

### Local Management (Direct Docker)

```bash
# Switch to local context
supactl config use-context local

# Manage local instances directly with Docker
supactl local add my-project
supactl list
supactl start my-project
```

**[Full supactl Documentation →](SUPACTL_README.md)**

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Examples](#quick-examples)
- [Features](#features)
- [Documentation](#documentation)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

`supactl` is a unified, Go-based command-line interface for managing self-hosted Supabase instances using a context-aware architecture inspired by kubectl.

### Context-Aware Architecture

Instead of separate command hierarchies, supactl uses contexts to switch between management modes:

- **Unified Commands**: Same commands work across all contexts
- **Seamless Switching**: Change contexts without learning new commands
- **kubectl-style UX**: Familiar commands like `get`, `describe`, and `config`

### Remote Management (SupaControl Server)

Connect to a SupaControl management server for centralized control:

- **Centralized Management**: Manage instances across multiple servers from one place
- **API-Driven**: Full REST API integration with the SupaControl server
- **Local Linking**: Connect local development directories to remote instances
- **Team Collaboration**: Share and manage instances with your team
- **Secure**: Encrypted credentials storage with proper permissions

### Local Management (Direct Docker)

Manage instances directly on your local machine with Docker:

- **No Server Required**: Works standalone, no additional infrastructure needed
- **Docker Integration**: Direct docker-compose management
- **Port Management**: Intelligent automatic port allocation
- **Quick Setup**: Automated Supabase project initialization with secure secrets
- **Full Control**: Direct access to your local instances

Both modes are available in a **single binary** that works cross-platform on Linux, macOS, and Windows.

## 🏗️ Architecture

### Context-Based Management

```bash
# View available contexts
supactl config get-contexts

# Switch between contexts
supactl config use-context local        # Local Docker mode
supactl config use-context production   # Remote SupaControl mode

# All commands work consistently across contexts
supactl list          # Lists instances in current context
supactl create myapp  # Creates remote instance (use 'local add' for local)
```

### kubectl-Style Commands

- `supactl get instances` - List instances in table format
- `supactl describe instance <name>` - Show detailed instance information
- `supactl config *` - Manage contexts (get-contexts, use-context, set-context, etc.)

## 📦 Installation

### Installing supactl

#### Using Installation Script (Recommended)

The installation script provides the most secure and user-friendly installation experience:

```bash
curl -sSL https://raw.githubusercontent.com/qubitquilt/supactl/main/scripts/install.sh | bash
```

**Benefits:**
- ✅ Automatically detects your platform (OS and architecture)
- ✅ Downloads and verifies SHA256 checksums
- ✅ Secure: Downloads to file first (never pipes remote content directly to tar)
- ✅ Handles permissions and installation to `/usr/local/bin` automatically
- ✅ Provides clear progress messages and error handling
- ✅ Works on Linux (amd64/arm64) and macOS (Intel/Apple Silicon)

#### Manual Installation (Advanced Users)

If you prefer manual installation with full control over each step:

```bash
# Linux (amd64)
curl -L "https://github.com/qubitquilt/supactl/releases/latest/download/supactl_Linux_x86_64.tar.gz" -o supactl.tar.gz
tar -xzf supactl.tar.gz
chmod +x supactl
sudo mv supactl /usr/local/bin/
rm supactl.tar.gz

# Linux (arm64)
curl -L "https://github.com/qubitquilt/supactl/releases/latest/download/supactl_Linux_arm64.tar.gz" -o supactl.tar.gz
tar -xzf supactl.tar.gz
chmod +x supactl
sudo mv supactl /usr/local/bin/
rm supactl.tar.gz

# macOS (Intel)
curl -L "https://github.com/qubitquilt/supactl/releases/latest/download/supactl_Darwin_x86_64.tar.gz" -o supactl.tar.gz
tar -xzf supactl.tar.gz
chmod +x supactl
sudo mv supactl /usr/local/bin/
rm supactl.tar.gz

# macOS (Apple Silicon)
curl -L "https://github.com/qubitquilt/supactl/releases/latest/download/supactl_Darwin_arm64.tar.gz" -o supactl.tar.gz
tar -xzf supactl.tar.gz
chmod +x supactl
sudo mv supactl /usr/local/bin/
rm supactl.tar.gz
```

#### Using Homebrew (macOS)
```bash
brew tap qubitquilt/homebrew-tap
brew install supactl
```

#### From Source
```bash
git clone https://github.com/qubitquilt/supactl.git
cd supactl
make build
sudo make install
```

## 🎨 Quick Examples

### Context Setup and Management

```bash
# Set up local context for Docker-based management
supactl config set-context local --provider=local
supactl config use-context local

# Set up remote context for SupaControl server
supactl config set-context production --provider=remote \
  --server=https://supacontrol.example.com --api-key=sk_...

# List all contexts
supactl config get-contexts
```

### Working with Instances

```bash
# Create remote instances (use 'local add' for local)
supactl create my-project
supactl create staging-environment

# List instances with kubectl-style output
supactl get instances

# Get detailed information
supactl describe instance my-project

# Manage lifecycle
supactl start my-project
supactl stop my-project
supactl restart my-project

# View logs
supactl logs my-project --lines 50
```

### Remote Mode with Project Linking

```bash
# Switch to remote context
supactl config use-context production

# Link local directory to remote instance
cd ~/my-project
supactl link
supactl status
```

### Local Mode with Docker

```bash
# Switch to local context
supactl config use-context local

# Create and manage local instances
supactl local add my-local-project
supactl list
supactl start my-local-project
```

## ✨ Features

### supactl Features

- **Context-Aware Architecture**: Unified commands across remote and local modes
- **kubectl-Style Commands**: Familiar commands like `get`, `describe`, and `config`
- **Authentication & Configuration**: Secure login with API key authentication
- **Instance Management**: Create, list, delete, start, stop, and restart instances
- **Local Project Linking**: Link development directories to remote instances (remote mode)
- **Status Monitoring**: View detailed instance information
- **Security**: Credentials stored with 600 permissions, no key echoing
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Single Binary**: No runtime dependencies

### Local Mode Features

- **Docker Integration**: Direct docker-compose management
- **Port Management**: Automatic unique port allocation
- **Secrets Generation**: Auto-generated secure passwords and JWT tokens
- **Project Isolation**: Each instance in separate Docker environment
- **Quick Setup**: Automated Supabase project initialization

## 📚 Documentation

### Detailed Documentation

- **[supactl Full Documentation](SUPACTL_README.md)** - Complete guide for the CLI tool
  - All commands and options
  - Context management
  - API endpoints reference
  - Troubleshooting guide

### API Documentation

supactl works with SupaControl servers that implement these endpoints:

- `GET /api/v1/auth/me` - Validate authentication
- `GET /api/v1/instances` - List instances
- `POST /api/v1/instances` - Create instance
- `GET /api/v1/instances/<name>` - Get instance details
- `DELETE /api/v1/instances/<name>` - Delete instance
- `POST /api/v1/instances/<name>/start` - Start instance
- `POST /api/v1/instances/<name>/stop` - Stop instance
- `POST /api/v1/instances/<name>/restart` - Restart instance
- `GET /api/v1/instances/<name>/logs` - Get logs

## 🛠 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/qubitquilt/supactl.git
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
│   ├── root.go         # Root command and context management
│   ├── create.go       # Create instance command
│   ├── delete.go       # Delete instance command
│   ├── list.go         # List instances command
│   ├── start.go        # Start instance command
│   ├── stop.go         # Stop instance command
│   ├── restart.go      # Restart instance command
│   ├── logs.go         # View logs command
│   ├── get.go          # kubectl-style get command
│   ├── describe.go     # kubectl-style describe command
│   ├── config.go       # Context management commands
│   ├── link.go         # Link directory command (remote mode)
│   ├── unlink.go       # Unlink directory command (remote mode)
│   ├── status.go       # Show linked project status (remote mode)
│   ├── local.go        # Local mode parent command
│   ├── local_add.go    # Local: create new instance
│   ├── local_list.go   # Local: list instances
│   ├── local_start.go  # Local: start instance
│   ├── local_stop.go   # Local: stop instance
│   └── local_remove.go # Local: remove instance
├── internal/
│   ├── api/            # API client implementation
│   ├── auth/           # Authentication & config management
│   ├── provider/       # Context-aware provider interface
│   ├── link/           # Local project linking
│   └── local/          # Local instance management
├── scripts/            # Installation and utility scripts
├── main.go
├── Makefile
└── README.md
```

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

MIT License - Copyright (c) 2025 Qubit Quilt

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the conditions in the MIT License.

See [LICENSE](LICENSE) for full details.

## 🙏 Acknowledgments

- Original development by Qubit Quilt
- Built on top of [Supabase](https://supabase.com/)
- CLI framework powered by [Cobra](https://github.com/spf13/cobra)

## 💬 Support & Community

- **Issues**: [GitHub Issues](https://github.com/qubitquilt/supactl/issues)
- **Discussions**: [GitHub Discussions](https://github.com/qubitquilt/supactl/discussions)
- **Documentation**: See [SUPACTL_README.md](SUPACTL_README.md) for detailed CLI docs

## 🗺️ Roadmap

### Completed
- [x] Context-aware multi-provider architecture (local Docker + remote SupaControl)
- [x] kubectl-style commands (get instances, describe instance, config management)
- [x] Local mode integration (direct Docker Compose management)
- [x] Multi-context support (multiple remote servers + local)
- [x] Secure local project linking (`.supacontrol/project` files)

### Planned
- [ ] Web UI for instance management
- [ ] Instance templates and presets
- [ ] Automated backups and restoration
- [ ] Multi-region support
- [ ] Instance monitoring and alerts
- [ ] Database migration tools
- [ ] Team collaboration features
- [ ] Bash/Zsh completion scripts
- [ ] Config file encryption at rest
- [ ] Instance health checking
- [ ] Bulk operations
- [ ] Declarative management (apply -f)


## ⚡ Performance

- **Fast**: Single binary with minimal overhead
- **Lightweight**: ~10MB binary size
- **Efficient**: Concurrent API requests where possible
- **Reliable**: Comprehensive error handling and retry logic

---

<p align="center">Made with ❤️ by Qubit Quilt</p>
