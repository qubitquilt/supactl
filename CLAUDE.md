# CLAUDE.md - AI Assistant Context

This document provides comprehensive context about the **supactl** project for AI assistants like Claude. It describes the project structure, architecture, conventions, and important implementation details.

## Project Overview

**supactl** is a command-line interface (CLI) tool for managing self-hosted Supabase instances via a central SupaControl management server. The repository also includes `supascale.sh`, a standalone bash script for direct local management.

### Core Purpose

- Enable centralized management of multiple Supabase instances across different servers
- Provide a user-friendly CLI with interactive prompts and clear error messages
- Support local development workflows by linking project directories to remote instances
- Maintain security through encrypted credential storage and API key authentication

## Repository Structure

```
supactl/
├── cmd/                      # Cobra command implementations
│   ├── root.go              # Root command and shared utilities
│   ├── login.go             # Authentication: login to SupaControl server
│   ├── logout.go            # Authentication: logout/clear credentials
│   ├── create.go            # Instance management: create new instance
│   ├── delete.go            # Instance management: delete instance
│   ├── list.go              # Instance management: list all instances
│   ├── start.go             # Lifecycle: start stopped instance
│   ├── stop.go              # Lifecycle: stop running instance
│   ├── restart.go           # Lifecycle: restart instance
│   ├── logs.go              # Debugging: view instance logs
│   ├── link.go              # Local dev: link directory to instance
│   ├── unlink.go            # Local dev: unlink directory
│   ├── status.go            # Local dev: show linked instance details
│   └── *_test.go            # Command tests
├── internal/                 # Internal packages (not for import)
│   ├── api/                 # API client for SupaControl server
│   │   ├── client.go        # HTTP client implementation
│   │   ├── types.go         # Request/response types
│   │   └── *_test.go        # API client tests
│   ├── auth/                # Authentication and config management
│   │   ├── config.go        # Config file I/O, credential storage
│   │   └── *_test.go        # Auth tests
│   ├── link/                # Local project linking
│   │   ├── link.go          # .supacontrol/project file management
│   │   └── *_test.go        # Link tests
│   └── testutil/            # Shared test utilities
├── scripts/                  # Installation and utility scripts
│   ├── install.sh           # One-line installation script
│   └── uninstall.sh         # Clean uninstallation script
├── .github/workflows/        # CI/CD pipelines
│   ├── test.yml             # Test workflow (multi-platform)
│   └── release.yml          # Release workflow (GoReleaser)
├── main.go                   # Application entry point
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Build automation
├── .goreleaser.yml          # Release configuration
├── .golangci.yml            # Linter configuration
├── README.md                 # User-facing documentation
├── SUPACTL_README.md        # Detailed CLI documentation
├── CONTRIBUTING.md          # Contributor guidelines
├── LICENSE                   # GPL v3 license
├── supascale.sh             # Standalone bash management script
└── CLAUDE.md                # This file

```

## Architecture

### Command Pattern

The CLI uses the **Cobra** framework for command structure:

1. **Root Command** (`cmd/root.go`): Base command with shared utilities
2. **Subcommands**: Each feature is a separate command file
3. **Shared Helper**: `getAPIClient()` creates authenticated API client

### API Client

**Location**: `internal/api/`

The API client (`client.go`) provides methods for all SupaControl server interactions:

- **Authentication**: `LoginTest()` validates API keys
- **Instance CRUD**: `CreateInstance()`, `GetInstance()`, `DeleteInstance()`, `ListInstances()`
- **Lifecycle**: `StartInstance()`, `StopInstance()`, `RestartInstance()`
- **Debugging**: `GetLogs()`

**Key Design Decisions**:
- Uses standard `net/http` package (no external HTTP libraries)
- 30-second timeout for all requests
- Centralized error handling via `handleErrorResponse()`
- Bearer token authentication in `Authorization` header

### Authentication & Configuration

**Location**: `internal/auth/config.go`

Configuration is stored in `~/.supacontrol/config.json` with **0600 permissions** (Unix) for security:

```json
{
  "server_url": "https://supacontrol.example.com",
  "api_key": "user-api-key"
}
```

**Cross-Platform Considerations**:
- Unix: Uses `$HOME` environment variable
- Windows: Uses `$USERPROFILE` or `$HOME`
- File permissions enforced only on Unix (Windows uses ACLs differently)

### Local Project Linking

**Location**: `internal/link/link.go`

Links local development directories to remote instances by creating `.supacontrol/project` file:

```
your-project/
└── .supacontrol/
    └── project        # Contains: "my-project-name"
```

**Features**:
- Automatically adds `.supacontrol/` to `.gitignore` if present
- Used by `status` command to show details of linked project
- Allows directory-based context instead of specifying project names repeatedly

## Code Conventions

### Naming

- **Commands**: Verb-based (login, create, delete, start, stop)
- **API Methods**: PascalCase with action prefix (CreateInstance, GetLogs)
- **Files**: Snake_case for multi-word files, lowercase for single words
- **Tests**: `*_test.go` files alongside implementation

### Error Handling

**Pattern Used Throughout**:
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

- Always print errors to `stderr` (not `stdout`)
- Prefix with "Error: " for clarity
- Exit with code 1 on failure
- API client returns user-friendly error messages (not raw HTTP codes)

### Testing

**Test Organization**:
- Unit tests alongside implementation (`*_test.go`)
- Table-driven tests for multiple scenarios
- Temporary directories for file I/O tests
- Mock HTTP servers for API client tests

**Platform-Specific Tests**:
```go
// Skip permission checks on Windows
if runtime.GOOS != "windows" {
    // Assert file permissions
}
```

**Windows Compatibility**:
- Set both `HOME` and `USERPROFILE` environment variables in tests
- Add small delays (`time.Sleep`) for file handle cleanup on Windows
- Skip Unix-specific assertions (file permissions)

### Project Validation

Project names must match regex: `^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`

**Valid**: `my-project`, `api-v2`, `test123`, `a`
**Invalid**: `MyProject`, `-start`, `end-`, `my_project`

## Dependencies

### Direct Dependencies

- **github.com/spf13/cobra** v1.8.0 - CLI framework
- **github.com/AlecAivazis/survey/v2** v2.3.7 - Interactive prompts

### Why These Choices?

- **Cobra**: Industry-standard CLI framework, excellent documentation, wide adoption
- **Survey**: Best-in-class terminal UI for interactive prompts, cross-platform
- **No HTTP library**: Standard library sufficient, reduces dependencies
- **No JSON library**: Standard `encoding/json` is performant and well-tested

## API Contract

### SupaControl Server Endpoints

The CLI expects these endpoints from the SupaControl server:

| Method | Endpoint | Purpose | Request | Response |
|--------|----------|---------|---------|----------|
| GET | `/api/v1/auth/me` | Validate authentication | - | `AuthResponse` |
| GET | `/api/v1/instances` | List instances | - | `ListInstancesResponse` |
| POST | `/api/v1/instances` | Create instance | `CreateInstanceRequest` | `Instance` |
| GET | `/api/v1/instances/:name` | Get instance details | - | `Instance` |
| DELETE | `/api/v1/instances/:name` | Delete instance | - | - |
| POST | `/api/v1/instances/:name/start` | Start instance | - | - |
| POST | `/api/v1/instances/:name/stop` | Stop instance | - | - |
| POST | `/api/v1/instances/:name/restart` | Restart instance | - | - |
| GET | `/api/v1/instances/:name/logs?lines=N` | Get logs | - | Plain text |

**Authentication**: All requests include `Authorization: Bearer <api_key>` header

### Type Definitions

See `internal/api/types.go` for complete type definitions:

- `Instance`: Represents a Supabase instance
- `AuthResponse`: Authentication validation response
- `ErrorResponse`: Standard error format
- `ListInstancesResponse`: Array of instances
- `CreateInstanceRequest`: Instance creation payload

## Build & Release

### Local Build

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run tests
make lint           # Run linter
```

### Release Process

**Automated via GitHub Actions**:

1. Tag a release: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions triggers `.github/workflows/release.yml`
3. GoReleaser builds binaries for all platforms
4. Creates GitHub release with:
   - Binaries (Linux, macOS, Windows for amd64/arm64)
   - Archives (tar.gz, zip)
   - Checksums
   - Auto-generated changelog

**Manual Release**:
```bash
goreleaser release --clean
```

## Testing Strategy

### Test Coverage

- **cmd/**: Command validation and flow tests
- **internal/api/**: Mock HTTP server tests for all endpoints
- **internal/auth/**: Config file I/O, permissions, cross-platform paths
- **internal/link/**: File creation, .gitignore updates

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Specific package
go test ./internal/api -v

# Platform-specific
GOOS=windows go test ./...  # Test Windows compatibility
```

### CI/CD

**GitHub Actions** (`.github/workflows/test.yml`):
- Tests on Ubuntu, macOS, Windows
- Tests with Go 1.21 and 1.22
- Runs linter (golangci-lint)
- Uploads coverage to Codecov
- Builds multi-platform binaries

## Important Implementation Details

### Security Considerations

1. **API Keys Never Logged**: Survey library used with `Password` type for input
2. **Config File Permissions**: 0600 on Unix systems (owner read/write only)
3. **No Secrets in Code**: All credentials come from config file or user input
4. **HTTPS Validation**: URLs must start with `http://` or `https://`

### Cross-Platform Compatibility

**File Paths**:
- Always use `filepath.Join()` for path construction
- Use `os.UserHomeDir()` for home directory (works on all platforms)

**Environment Variables**:
- Unix: `$HOME`
- Windows: `$USERPROFILE` or `$HOME`
- Always set both in tests

**Shell Commands**:
- GitHub Actions: Force `bash` shell on Windows to avoid PowerShell parsing issues
- Race detector: Disabled on Windows (not fully supported)

### User Experience

**Interactive Prompts**:
- Login: Password prompt (no echo) for API key
- Delete: Confirmation prompt before destructive actions
- Link: Selection menu for choosing from available instances

**Output Formatting**:
- List command: Tab-formatted table with aligned columns
- Error messages: Always prefixed with "Error: "
- Success messages: Clear confirmation of action taken
- Status command: Hierarchical display with indentation

## Common Development Tasks

### Adding a New Command

1. Create `cmd/newcommand.go`:
```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new <args>",
    Short: "Brief description",
    Long:  `Detailed description`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        client := getAPIClient()
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

2. Add API method to `internal/api/client.go` if needed
3. Add tests to `cmd/newcommand_test.go`
4. Update documentation in `SUPACTL_README.md`

### Adding a New API Endpoint

1. Add method to `internal/api/client.go`
2. Add types to `internal/api/types.go` if needed
3. Add tests to `internal/api/client_test.go`
4. Use in command implementation

### Modifying Configuration

1. Update `Config` struct in `internal/auth/config.go`
2. Update `SaveConfig()` and `LoadConfig()` functions
3. Update tests in `internal/auth/config_test.go`
4. Consider migration path for existing users

## Troubleshooting

### Common Issues

**"You are not logged in"**:
- User needs to run `supactl login <server_url>` first
- Config file missing or invalid at `~/.supacontrol/config.json`

**"No project linked"**:
- User needs to run `supactl link` in project directory
- `.supacontrol/project` file missing or empty

**Tests failing on Windows**:
- Ensure USERPROFILE is set in tests
- Skip Unix-specific assertions (permissions)
- Add file handle cleanup delays if needed

**Port conflicts** (supascale.sh):
- Each project uses base port + 1000 increments
- Check for existing Docker containers: `docker ps`

## Philosophy & Design Principles

1. **User-Centric**: Clear error messages, helpful prompts, good defaults
2. **Secure by Default**: Credentials never logged, files have restrictive permissions
3. **Cross-Platform**: Works identically on Linux, macOS, and Windows
4. **Self-Contained**: Single binary, no runtime dependencies
5. **API-Driven**: All business logic lives on server, CLI is thin client
6. **Unix Philosophy**: Each command does one thing well, composable
7. **Fail Fast**: Exit on errors with clear messages, don't continue in invalid state

## Future Enhancements

- [ ] Bash/Zsh completion scripts
- [ ] Config file encryption at rest
- [ ] Support for multiple server profiles
- [ ] Instance health checking
- [ ] Bulk operations (create/delete multiple instances)
- [ ] Instance import/export
- [ ] Plugin system for extensions
- [ ] TUI (terminal UI) mode

## License

GPL v3 - See LICENSE file for full text. Copyright (c) 2025 Frog Byte, LLC.

---

**Last Updated**: 2025-11-07
**Version**: 1.0.0
**Maintained By**: Frog Byte, LLC
