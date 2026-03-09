# Muster CLI

> A beautiful terminal user interface (TUI) for Muster - a team productivity tool for daily standups, attendance tracking, and leave management.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

---

## Features

- 📝 **Standups**: Submit and view team standups
- 📊 **Attendance**: Check in/out and view team status
- 🏖️ **Leave Management**: Request, review, and cancel leaves
- 👥 **Team Management**: Invite members, manage roles
- 🔐 **Authentication**: Login, registration, and org invitations
- 🎨 **Beautiful TUI**: Built with Bubble Tea and Lip Gloss

## Installation

### Homebrew (macOS & Linux)

```bash
brew install beamlabco/tap/muster
```

### Build from source

```bash
git clone https://github.com/beamlabco/muster.git
cd muster
make build
./bin/muster
```

## Usage

Launch the CLI:

```bash
muster
```

### Commands

Type `/` to see available commands:

| Command | Description |
|---------|-------------|
| `/login` | Sign in to your account |
| `/signup` | Create a new account |
| `/join` | Join an organization via invitation |
| `/standup` | Submit your daily standup |
| `/standup today` | View team standups for today |
| `/standup history` | View your standup history |
| `/checkin` | Check in for the day |
| `/checkout` | Check out for the day |
| `/attendance today` | View team attendance for today |
| `/attendance history` | View your attendance history |
| `/leave` | Request a leave |
| `/leave list` | View team leaves |
| `/leave review` | Review pending leaves (primary only) |
| `/leave cancel` | Cancel a pending leave |
| `/team` | View team members and roles |
| `/role` | Update a user's role (primary only) |
| `/invite` | Invite a team member |
| `/whoami` | Show current user info |
| `/help` | Show available commands |
| `/quit` | Exit the application |

### Navigation

- `Tab`: Accept autocomplete suggestion
- `↑` / `↓`: Navigate suggestions
- `Enter`: Execute command
- `Esc`: Go back / cancel
- `Ctrl+C`: Quit

### Configuration

Config is stored at `~/.muster/config.yaml`:

```yaml
user:
  id: 1
  email: user@example.com
  name: John Doe
  role: primary

organization:
  id: 1
  name: Acme Corp
  status: active

auth:
  token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

api:
  base_url: https://api.muster.stagify.xyz
```

## Development

### Make commands

```bash
make build         # Build the CLI binary
make run           # Run the CLI directly
make deps          # Install dependencies
make install       # Install binary to /usr/local/bin
make clean         # Clean build artifacts
make test          # Run tests
make test-coverage # Run tests with coverage
make lint          # Run linters (requires golangci-lint)
make fmt           # Format code
make build-all     # Build for all platforms
```

### Releasing

Releases are automated via GitHub Actions and [GoReleaser](https://goreleaser.com/). Pushing a version tag triggers a build that publishes binaries to GitHub Releases and updates the Homebrew tap.

```bash
# Create and push a release tag
git tag v0.2.0
git push origin v0.2.0
```

To test a release locally without publishing:

```bash
goreleaser release --snapshot --clean
```

Binaries are built for:
- macOS (Intel & Apple Silicon)
- Linux (amd64 & arm64)
- Windows (amd64)

## License

See the main project LICENSE file.
