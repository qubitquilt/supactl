# supactl

A command-line interface for managing self-hosted Supabase instances via a central SupaControl server.

## Overview

`supactl` (Supa-Control CLI) provides a simple, intuitive, and powerful way to manage the lifecycle of self-hosted Supabase instances from your terminal. It connects to a SupaControl management server and allows you to create, list, delete, and manage Supabase instances seamlessly.

## Features

- **Authentication & Configuration**: Secure login to your SupaControl server with API key authentication
- **Instance Management**: Create, list, and delete Supabase instances
- **Local Project Linking**: Link your local development directory to a remote instance
- **Secure by Default**: Credentials stored with 600 permissions, API keys never echoed
- **Single Binary**: Self-contained executable with no external dependencies

## Installation

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd supactl

# Build the binary
go build -o supactl

# Move to your PATH (optional)
sudo mv supactl /usr/local/bin/
```

### Pre-built Binaries

Download the latest release for your platform from the releases page.

## Quick Start

### 1. Login to Your SupaControl Server

```bash
supactl login https://your-supacontrol-server.com
```

You'll be prompted to enter your API key (obtain this from your SupaControl dashboard).

### 2. Create a New Instance

```bash
supactl create my-project
```

### 3. List Your Instances

```bash
supactl list
```

### 4. Link Your Local Directory

```bash
cd /path/to/your/project
supactl link
```

Select your project from the interactive list.

### 5. Check Status

```bash
supactl status
```

## Commands

### Authentication

#### `supactl login <server_url>`

Login to your SupaControl server.

```bash
supactl login https://supacontrol.example.com
```

- Prompts for API key (no echo for security)
- Validates credentials before saving
- Stores config in `~/.supacontrol/config.json` with 600 permissions

#### `supactl logout`

Logout from your SupaControl server.

```bash
supactl logout
```

Removes stored credentials.

### Instance Management

#### `supactl create <project-name>`

Create a new Supabase instance.

```bash
supactl create my-new-project
```

**Project name requirements:**
- Lowercase letters and numbers
- May contain hyphens
- Must start and end with an alphanumeric character

Example valid names: `my-project`, `api-v2`, `test123`

#### `supactl list`

List all your Supabase instances.

```bash
supactl list
```

Displays a table with:
- Project name
- Status
- Studio URL

#### `supactl delete <project-name>`

Delete a Supabase instance.

```bash
supactl delete my-project
```

**Warning**: This action is irreversible. You'll be asked to confirm before deletion.

### Local Development Integration

#### `supactl link`

Link the current directory to a remote instance.

```bash
cd /path/to/your/project
supactl link
```

- Presents an interactive selector with your available instances
- Creates `.supacontrol/project` file in the current directory
- Automatically adds `.supacontrol/` to `.gitignore` if present

#### `supactl unlink`

Unlink the current directory from its remote instance.

```bash
supactl unlink
```

Removes the `.supacontrol/project` file.

#### `supactl status`

Show details about the linked project.

```bash
supactl status
```

Displays:
- Project name
- Status
- Studio URL
- API URL
- Kong URL
- Anon key (if available)
- Creation date

**Note**: Requires the directory to be linked first.

### Help & Information

#### `supactl --help`

Show help for all commands.

```bash
supactl --help
```

#### `supactl <command> --help`

Show detailed help for a specific command.

```bash
supactl create --help
```

#### `supactl --version`

Display the CLI version.

```bash
supactl --version
```

## Configuration

### Authentication Configuration

Credentials are stored in `~/.supacontrol/config.json`:

```json
{
  "server_url": "https://supacontrol.example.com",
  "api_key": "your-api-key"
}
```

**Security**: This file has 600 permissions (read/write for user only).

### Project Linking

When you link a directory, a `.supacontrol/project` file is created:

```
your-project/
└── .supacontrol/
    └── project
```

This file contains the name of the linked project.

## API Endpoints

`supactl` communicates with the following SupaControl API endpoints:

- `GET /api/v1/auth/me` - Validate authentication
- `GET /api/v1/instances` - List instances
- `POST /api/v1/instances` - Create instance
- `GET /api/v1/instances/<name>` - Get instance details
- `DELETE /api/v1/instances/<name>` - Delete instance

## Examples

### Complete Workflow

```bash
# Login
supactl login https://supacontrol.example.com
# Enter your API key when prompted

# Create a new project
supactl create my-awesome-app

# List all projects
supactl list

# Link to local directory
cd ~/projects/my-awesome-app
supactl link
# Select "my-awesome-app" from the list

# Check project status
supactl status

# When done, delete the project
supactl delete my-awesome-app
# Confirm deletion
```

### Working with Multiple Instances

```bash
# Create multiple instances
supactl create dev-environment
supactl create staging-environment
supactl create production-environment

# List them all
supactl list

# Link different directories to different instances
cd ~/projects/dev && supactl link    # Select dev-environment
cd ~/projects/staging && supactl link # Select staging-environment
cd ~/projects/prod && supactl link    # Select production-environment
```

## Troubleshooting

### "You are not logged in" Error

```bash
Error: You are not logged in. Please run 'supactl login <server_url>' first.
```

**Solution**: Run `supactl login <server_url>` to authenticate.

### "No project linked" Error

```bash
Error: No project linked. Run 'supactl link' to get started.
```

**Solution**: Run `supactl link` in your project directory.

### "Authentication failed" Error

**Possible causes**:
- Invalid API key
- Incorrect server URL
- Server is unreachable

**Solution**:
1. Verify your server URL is correct
2. Obtain a fresh API key from your SupaControl dashboard
3. Check network connectivity

### Invalid Project Name

```bash
Error: Project name 'MyProject' is invalid.
Name must be lowercase, alphanumeric, and may contain hyphens.
```

**Solution**: Use only lowercase letters, numbers, and hyphens. Example: `my-project`

## Development

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd supactl

# Install dependencies
go mod download

# Build
go build -o supactl

# Run
./supactl --help
```

### Project Structure

```
supactl/
├── cmd/                  # Cobra command definitions
│   ├── root.go          # Root command
│   ├── login.go         # Login command
│   ├── logout.go        # Logout command
│   ├── create.go        # Create instance command
│   ├── list.go          # List instances command
│   ├── delete.go        # Delete instance command
│   ├── link.go          # Link directory command
│   ├── unlink.go        # Unlink directory command
│   └── status.go        # Show status command
├── internal/
│   ├── api/             # API client
│   │   ├── client.go   # HTTP client implementation
│   │   └── types.go    # API types
│   ├── auth/            # Authentication
│   │   └── config.go   # Config management
│   └── link/            # Local linking
│       └── link.go     # Link file management
├── main.go              # Application entry point
└── go.mod               # Go module definition
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- Go standard library

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - Copyright (c) 2025 Qubit Quilt

## Acknowledgments

- Built for [SupaControl](https://github.com/qubitquilt/supacontrol)
- Powered by [Supabase](https://supabase.com/)
- CLI framework by [Cobra](https://github.com/spf13/cobra)

## Support

For issues, questions, or contributions:
- GitHub Issues: [repository-url/issues]
- Documentation: [docs-url]

## Constitution & Philosophy

**User-Centric**: The CLI is intuitive and designed for humans. Clear commands, helpful error messages, and interactive prompts make it accessible.

**Secure by Default**: API keys are never echoed, logged, or stored in world-readable files. All credentials use 600 permissions.

**Self-Contained**: Single binary with no external runtime dependencies.

**API-Driven**: The CLI is a client for the SupaControl server. All business logic lives on the server.

**Unix Philosophy**: Each command does one thing well. Commands can be composed for complex workflows.
